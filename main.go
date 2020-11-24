package main

import (
	"rpiWater/gpioControl"
	"rpiWater/socket"
	"rpiWater/sysSignal"
)

func main() {

	go gpioControl.Run()

	go sockets.Run()

	sysSignal.Run()
}
