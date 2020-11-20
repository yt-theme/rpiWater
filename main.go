package main

import (
    "rpiWater/socket"
    "rpiWater/gpioControl"
)

func main () {
    
    go gpioControl.Run()
    
    sockets.Run()
}