GO WASM Elevator Saga
===================
The [elevator programming game](https://github.com/magwo/elevatorsaga), modified to accept GO WASM input:
- The user submits his go code 
- The code is sent to a go API server
- The API server builds the wasm binary in a docker container
- The wasm binary is returned to the browser and loaded
- Frontend JS calls the wasm code and runs the solution

**Only tested on Chrome v71+ !**

TODO: 
- try out tinygo

[Play it now!](https://didil.github.io/gowasm-elevatorsaga/)

![Image of Elevator Saga in browser](https://raw.githubusercontent.com/didil/gowasm-elevatorsaga/master/images/screenshot.png)