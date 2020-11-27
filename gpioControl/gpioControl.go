package gpioControl

import (
	"fmt"
	"rpiWater/public"
	"time"
)

var IsStarted = false

func init() {

	// --------------------------------------------------------------------
	// init gpio
	initGpioRet, err := ExecPython([]string{"gpio.py", "INIT"})
	if err != nil {
		fmt.Println("initGpio Err =>", err)
	}
	fmt.Println(initGpioRet)
}

/*
   todo: 液位低传感器
*/

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
   sensor1(液位满传感器) 状态检查
*/
func checkSensor1() {
	public.S1_isChecking = true
	var is_fullToLow = false // 代表从满到空的过程
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
			gpioRet18, err := ExecPython([]string{"gpio.py", "18", "IN", "READ"}) // S1(液位传感器)
			if err != nil {
				fmt.Println("gpioRet sensor1 Err =>", err)
			}
			// fmt.Println("S1 Value ======>", gpioRet)

			// 满则关闭泵1(向小缸抽水) 延时开户泵2(从小缸排水)
			// fmt.Println("do close M1 ===============================>", gpioRet, len(gpioRet), gpioRet == "1")

			if gpioRet18 == "0" {
				is_fullToLow = true
				_, err := ExecPython([]string{"gpio.py", "12", "OUT", "LOW"})
				if err != nil {
					fmt.Println("gpioRet pin12LOW Err =>", err)
				}
				// fmt.Println(gpioRet1)

				time.Sleep(100000000)

				_, err = ExecPython([]string{"gpio.py", "16", "OUT", "HIGH"})
				if err != nil {
					fmt.Println("gpioRet pin16HIGH Err =>", err)
				}
				// fmt.Println(gpioRet)
			} else if gpioRet18 == "1" {
				// 如果S1高电平 && S1高电平 && 是从满到空的过程 则关闭进程
				gpioRet22 := checkSensor2()
				if gpioRet18 == "1" && gpioRet22 == "1" && is_fullToLow == true {
					fmt.Println("do stop 16 =====>")
					ExecPython([]string{"gpio.py", "16", "OUT", "LOW"})

					goto endProcess
				}
			}
			// --------------------------------------------------------
			break
		}
	}
end:
	public.S1_isChecking = false
	return
endProcess:
	// public.Chan_stop <- 1
	public.Chan_sendMsg <- "cur process completed."
	IsStarted = false

	return
}

/*
   sensor2(液位空传感器) 状态检查
*/
func checkSensor2() string {
	gpioRet, err := ExecPython([]string{"gpio.py", "22", "IN", "READ"}) // S1(液位传感器)
	if err != nil {
		fmt.Println("gpioRet sensor2 Err =>", err)
	}
	fmt.Println("S2 Value ======>", gpioRet)

	return gpioRet
}

// -----------------------------------------------------------------------------
/*
   0 close
   1 open
*/
// 总体进程操作
func handleProcess(opa int) {
	// 进程结束操作
	if opa == 0 {
		if IsStarted == true {
			public.Chan_S1_stopCheck <- 1 // sensor1停止检测

			// ------------------------------------------------------------
			ExecPython([]string{"gpio.py", "12", "OUT", "LOW"})
			ExecPython([]string{"gpio.py", "16", "OUT", "LOW"})

			// gpioRet, err := ExecPython([]string{"gpio.py", "CLEANUP", "", ""})
			// if err != nil {
			// 	fmt.Println("gpioRet stopCheck Err =>", err)
			// }
			// fmt.Println(gpioRet)
			public.Chan_sendMsg <- "process stop"
			// ------------------------------------------------------------

		} else {
			ExecPython([]string{"gpio.py", "12", "OUT", "LOW"})
			ExecPython([]string{"gpio.py", "16", "OUT", "LOW"})
			public.Chan_sendMsg <- "process already stopped"
		}

		IsStarted = false
	}
	// 进程开始
	if opa == 1 {
		fmt.Println("do process start =>")
		if IsStarted == false {
			public.Chan_sendMsg <- "process start"
			ExecPython([]string{"gpio.py", "12", "OUT", "HIGH"})
			go checkSensor1() // sensor1开始检测
		} else {
			public.Chan_sendMsg <- "process already started"
		}

		IsStarted = true
	}
}

// 继电器1 (向小缸抽水)
func handleR1(opa int) {
	// 断开
	if opa == 0 {
		// ------------------------------------------------------------
		gpioRet, err := ExecPython([]string{"gpio.py", "12", "OUT", "LOW"})
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
		gpioRet, err := ExecPython([]string{"gpio.py", "12", "OUT", "HIGH"})
		if err != nil {
			fmt.Println("gpioRet handleR1 Err =>", err)
		}
		fmt.Println(gpioRet)
		public.Chan_sendMsg <- "R1 on"
		// ------------------------------------------------------------
	}
}

// 继电器2 (从小缸排水)
func handleR2(opa int) {
	if opa == 0 {
		// ------------------------------------------------------------
		gpioRet, err := ExecPython([]string{"gpio.py", "16", "OUT", "LOW"})
		if err != nil {
			fmt.Println("gpioRet handleR2 Err =>", err)
		}
		fmt.Println(gpioRet)
		public.Chan_sendMsg <- "R2 close"
		// ------------------------------------------------------------
	}
	if opa == 1 {
		// ------------------------------------------------------------
		gpioRet, err := ExecPython([]string{"gpio.py", "16", "OUT", "HIGH"})
		if err != nil {
			fmt.Println("gpioRet handleR2 Err =>", err)
		}
		fmt.Println(gpioRet)
		public.Chan_sendMsg <- "R2 on"
		// ------------------------------------------------------------
	}
}
