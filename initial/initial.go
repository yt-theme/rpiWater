package initial

import (
	"fmt"
	"rpiWater/gpioControl"
)

func Run(callback func()) {
	// low gpio
	gpioRet1, err := gpioControl.ExecPython([]string{"gpio.py", "12", "OUT", "LOW"})
	if err != nil {
		fmt.Println("gpio 12 low Err =>", err)
	}
	fmt.Println("gpio 12 low ======>", gpioRet1)

	gpioRet2, err := gpioControl.ExecPython([]string{"gpio.py", "16", "OUT", "LOW"})
	if err != nil {
		fmt.Println("gpio 16 low Err =>", err)
	}
	fmt.Println("gpio 16 low ======>", gpioRet2)

	callback()
}
