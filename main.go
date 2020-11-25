package main

import (
	"rpiWater/gpioControl"
	"rpiWater/initial"
	"rpiWater/socket"
	"rpiWater/sysSignal"
)

func main() {

	initial.Run(func() {
		go gpioControl.Run()

		go sockets.Run()

		sysSignal.Run()
	})
}
