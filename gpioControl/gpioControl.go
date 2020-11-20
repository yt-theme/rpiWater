package gpioControl

import (
	"fmt"
	"rpiWater/public"
	"time"
)

var isStarted = false

func init() {

	// --------------------------------------------------------------------
	// init gpio
	initGpioRet, err := execPython([]string{"gpio.py", "INIT"})
	if err != nil {
		fmt.Println("initGpio Err =>", err)
	}
	fmt.Println(initGpioRet)
}

func Run() {
	for {
		select {
		case <-public.Chan_start:
			handleProcess(1)

			break
		case <-public.Chan_stop:
			handleProcess(0)

			break

		// ----------------------------
		case <-public.Chan_R1_isActi:
			handleR1(1)

			break
		case <-public.Chan_R1_isInActi:
			handleR1(0)

			break

		case <-public.Chan_R2_isActi:
			handleR2(1)

			break
		case <-public.Chan_R2_isInActi:
			handleR2(0)

			break
		}
	}

}

// ========================================================================
//             handle
// ========================================================================
/*
   sensor1 状态检查
*/
func checkSensor1() {

	public.S1_isChecking = true
	for {
		select {
		case <-public.Chan_S1_stopCheck:
			goto end
			break
		default:
			fmt.Println("check sensor =>")
			time.Sleep(100000000)

			// --------------------------------------------------------
			// check func

			gpioRet, err := execPython([]string{"gpio.py", "18", "IN", "READ"}) // S1(液位传感器)
			if err != nil {
				fmt.Println("gpioRet sensor Err =>", err)
			}
			fmt.Println("S1 Value ======>", gpioRet)

			// 满则关闭泵1 延时开户泵2
			fmt.Println("do close M ===============================>", gpioRet, len(gpioRet), gpioRet == "1")

			if gpioRet == "0" {
				gpioRet1, err := execPython([]string{"gpio.py", "12", "OUT", "LOW"})
				if err != nil {
					fmt.Println("gpioRet pin12LOW Err =>", err)
				}
				fmt.Println(gpioRet1)

				time.Sleep(100000000)

				gpioRet, err = execPython([]string{"gpio.py", "16", "OUT", "HIGH"})
				if err != nil {
					fmt.Println("gpioRet pin16HIGH Err =>", err)
				}
				fmt.Println(gpioRet)
			} else {

			}
			// --------------------------------------------------------

			break
		}
	}

end:
	public.S1_isChecking = false
	return
}

/*
   0 close
   1 open
*/
// 总体进程操作
func handleProcess(opa int) {
	// 进程结束操作
	if opa == 0 {
		if isStarted == true {
			public.Chan_S1_stopCheck <- 1 // sensor1停止检测

			// ------------------------------------------------------------
			execPython([]string{"gpio.py", "12", "OUT", "LOW"})
			execPython([]string{"gpio.py", "16", "OUT", "LOW"})

			gpioRet, err := execPython([]string{"gpio.py", "CLEANUP", "", ""})
			if err != nil {
				fmt.Println("gpioRet stopCheck Err =>", err)
			}
			fmt.Println(gpioRet)
			public.Chan_sendMsg <- "process stop"
			// ------------------------------------------------------------

		} else {
			public.Chan_sendMsg <- "process already stopped"
		}

		isStarted = false
	}
	// 进程开始
	if opa == 1 {
		if isStarted == false {
			public.Chan_sendMsg <- "process start"
			execPython([]string{"gpio.py", "12", "OUT", "HIGH"})
			go checkSensor1() // sensor1开始检测
		} else {
			public.Chan_sendMsg <- "process already started"
		}

		isStarted = true
	}
}

// 继电器1
func handleR1(opa int) {
	// 断开
	if opa == 0 {
		// ------------------------------------------------------------
		gpioRet, err := execPython([]string{"gpio.py", "12", "OUT", "LOW"})
		if err != nil {
			fmt.Println("gpioRet handleR1 Err =>", err)
		}
		fmt.Println(gpioRet)
		public.Chan_sendMsg <- "R1 close"
		// ------------------------------------------------------------
	}
	// 接通
	if opa == 1 {
		// ------------------------------------------------------------
		fmt.Println("OPA START ===============>")
		gpioRet, err := execPython([]string{"gpio.py", "12", "OUT", "HIGH"})
		if err != nil {
			fmt.Println("gpioRet handleR1 Err =>", err)
		}
		fmt.Println(gpioRet)
		public.Chan_sendMsg <- "R1 on"
		// ------------------------------------------------------------
	}
}

// 继电器2
func handleR2(opa int) {
	if opa == 0 {
		// ------------------------------------------------------------
		gpioRet, err := execPython([]string{"gpio.py", "16", "OUT", "LOW"})
		if err != nil {
			fmt.Println("gpioRet handleR2 Err =>", err)
		}
		fmt.Println(gpioRet)
		public.Chan_sendMsg <- "R2 close"
		// ------------------------------------------------------------
	}
	if opa == 1 {
		// ------------------------------------------------------------
		gpioRet, err := execPython([]string{"gpio.py", "16", "OUT", "HIGH"})
		if err != nil {
			fmt.Println("gpioRet handleR2 Err =>", err)
		}
		fmt.Println(gpioRet)
		public.Chan_sendMsg <- "R2 on"
		// ------------------------------------------------------------
	}
}
