package api

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

// Start the api server
func Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/compile", handleCompile)
	return http.ListenAndServe(":3000", mux)
}

func handleCompile(w http.ResponseWriter, r *http.Request) {
	dir, err := ioutil.TempDir("/tmp/gowasmbuilder", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mainF, err := os.Open("./gowasm/main.go")
	if err != nil {
		handleErr(w, err)
		return
	}

	mainFCopy, err := os.Create(filepath.Join(dir, "main.go"))
	if err != nil {
		handleErr(w, err)
		return
	}

	_, err = io.Copy(mainFCopy, mainF)
	if err != nil {
		handleErr(w, err)
		return
	}

	inputF, err := os.Create(filepath.Join(dir, "input.go"))
	if err != nil {
		handleErr(w, err)
		return
	}

	_, err = io.Copy(inputF, r.Body)
	if err != nil {
		handleErr(w, err)
		return
	}
}

func handleErr(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
