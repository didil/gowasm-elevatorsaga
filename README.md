GO WASM Elevator Saga
===================
The [elevator programming game](https://github.com/magwo/elevatorsaga), modified to accept GO WASM input:
- The user submits his go code 
- A POST request is sent to a REST API endpoint with the user code + a hash
- The reverse proxy checks the cache, in case of a cache HIT the next step is skipped 
- The Go API server builds the wasm binary in a docker container
- The wasm binary is returned to the browser and loaded
- Frontend JS calls the wasm code and runs the solution

**Only tested on Chrome v71+ !**

TODO: 
- try out tinygo

[Play it now!](https://didil.github.io/gowasm-elevatorsaga/)

![Image of Elevator Saga in browser](https://raw.githubusercontent.com/didil/gowasm-elevatorsaga/master/images/screenshot.png)