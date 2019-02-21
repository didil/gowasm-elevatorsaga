package main

import (
	"syscall/js"
)

// Init corresponds to the init function from http://play.elevatorsaga.com/documentation.html#docs
func Init(args []js.Value) {
	elevators := args[0]
	//	floors := args[1]

	elevator := elevators.Index(0) // Let's use the first elevator

	// Whenever the elevator is idle (has no more queued destinations) ...
	idleCb := js.NewCallback(func(args []js.Value) {
		elevator.Call("goToFloor", 0)
		elevator.Call("goToFloor", 1)
	})

	// Attach callback
	elevator.Call("on", "idle", idleCb)
}

// Update corresponds to the update function from http://play.elevatorsaga.com/documentation.html#docs
func Update(args []js.Value) {
	// We normally don't need to do anything here
}
