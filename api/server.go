package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/rs/cors"
)

// Start the api server
func Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/compile", handleCompile)

	// cors
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})
	router := c.Handler(mux)

	fmt.Println("Listening on port 3000")
	return http.ListenAndServe(":3000", router)
}

type compileRequest struct {
	Compiler string `json:"compiler,omitempty"`
	Input    string `json:"input,omitempty"`
}

func handleCompile(w http.ResponseWriter, r *http.Request) {
	cRequest := &compileRequest{}
	err := json.NewDecoder(r.Body).Decode(cRequest)
	if err != nil {
		handleErr(w, err)
		return
	}

	dir, err := createSourceFiles(cRequest)
	if err != nil {
		handleErr(w, err)
		return
	}

	err = runCompile(dir)
	if err != nil {
		handleErr(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/wasm")

	wasmF, err := os.Open(filepath.Join(dir, "app.wasm"))
	if err != nil {
		handleErr(w, err)
		return
	}
	defer wasmF.Close()

	io.Copy(w, wasmF)
}

func handleErr(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func createSourceFiles(cRequest *compileRequest) (string, error) {
	dir, err := ioutil.TempDir("/tmp/", "gowasmbuilder-")
	if err != nil {
		return "", err
	}

	fmt.Println("Saving data to ", dir)

	mainF, err := os.Open("./gowasm/main.go")
	if err != nil {
		return "", err
	}
	defer mainF.Close()

	mainFCopy, err := os.Create(filepath.Join(dir, "main.go"))
	if err != nil {
		return "", err
	}
	defer mainFCopy.Close()

	_, err = io.Copy(mainFCopy, mainF)
	if err != nil {
		return "", err
	}

	inputF, err := os.Create(filepath.Join(dir, "input.go"))
	if err != nil {
		return "", err
	}
	defer inputF.Close()

	_, err = inputF.WriteString(cRequest.Input)
	if err != nil {
		return "", err
	}

	return dir, nil
}

func runCompile(dir string) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "didil/gowasmbuilder",
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			mount.Mount{Type: mount.TypeBind, Source: dir, Target: "/app/"},
		},
	}, nil, "")
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case status := <-statusCh:
		if status.StatusCode != 0 {
			out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
			if err != nil {
				return err
			}

			logs, err := ioutil.ReadAll(out)
			if err != nil {
				return err
			}

			return fmt.Errorf("Container exited with code %d\nlogs: %v", status.StatusCode, string(logs))
		}
	}
	return nil
}
