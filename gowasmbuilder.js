window.GoWasmBuilder = {
    exitWasm: null,
    codeObj: null,
    mod: null,
    inst: null,
    apiRoot: null,
    async init(bytes) {
        GoWasmBuilder.go = new Go();

        let result = await WebAssembly.instantiate(bytes, GoWasmBuilder.go.importObject);
        GoWasmBuilder.mod = result.module;
        GoWasmBuilder.inst = result.instance;
    },
    run() {
        GoWasmBuilder.go.run(GoWasmBuilder.inst)
    },
    async getCodeObjFromCode(code) {
        let resp = await fetch(GoWasmBuilder.apiRoot + "/api/v1/compile", {
            method: 'POST',
            headers: {
                'Accept': 'application/json, application/xml, text/plain, text/html, *.*',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ compiler: "go1.11.5", input: code, })
        })

        if (!resp.ok){
            let text = await resp.text()
            throw new Error(text)
        }

        let bytes = await resp.arrayBuffer()

        if (GoWasmBuilder.exitWasm) {
            GoWasmBuilder.exitWasm()
            GoWasmBuilder.exitWasm = null
        }

        await GoWasmBuilder.init(bytes)

        GoWasmBuilder.run()
        let codeObj = GoWasmBuilder.codeObj
        if (typeof codeObj.init !== "function") {
            throw new Error("Code must contain an init function");
        }
        if (typeof codeObj.update !== "function") {
            throw new Error("Code must contain an update function");
        }

        return codeObj
    }
}

if (window.location.href.includes("localhost")){
    GoWasmBuilder.apiRoot = "http://localhost:3000";
}
else {
    GoWasmBuilder.apiRoot = "https://gowasm-elevatorsaga.leclouddev.com";
}