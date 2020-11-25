package sysSignal

import (
	"fmt"
	"os"
	"os/signal"
	"rpiWater/initial"
	"rpiWater/public"
	"syscall"
	"time"
)

func Run() {
	signal.Notify(
		public.SysSignal_chan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGUSR1,
		syscall.SIGUSR2)

	for {
		s := <-public.SysSignal_chan
		fmt.Println("sys signal:", s)
		go func() {
			public.Chan_stop <- 1 // stop process action

			time.Sleep(100000000)

			initial.Run(func() {
				os.Exit(1)
			})
		}()
	}
}
