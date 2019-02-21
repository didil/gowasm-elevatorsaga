window.GoWasmBuilder = {
    exitWasm: null,
    codeObj: null,
    mod: null,
    inst: null,
    apiRoot: null,
    async init(bytes) {
        // load the go module
        GoWasmBuilder.go = new Go();

        let result = await WebAssembly.instantiate(bytes, GoWasmBuilder.go.importObject);
        GoWasmBuilder.mod = result.module;
        GoWasmBuilder.inst = result.instance;
    },
    run() {
        // run the go module
        GoWasmBuilder.go.run(GoWasmBuilder.inst)
    },
    async getCodeObjFromCode(code) {
        // build json input
        let json = JSON.stringify({ compiler: "go1.11.5", input: code, })
        // hash json input
        let hash = SparkMD5.hash(json); 
        // perform POST request
        let resp = await fetch(GoWasmBuilder.apiRoot + "/api/v1/compile", {
            method: 'POST',
            headers: {
                'Accept': 'application/json, application/xml, text/plain, text/html, *.*',
                'Accept-Encoding': 'gzip',
                'Content-Type': 'application/json',
                'Code-Hash': hash,
            },
            body: json
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

        // init module and run
        await GoWasmBuilder.init(bytes)
        GoWasmBuilder.run()

        // get result object
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