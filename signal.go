package stop

import (
	"os"
	"os/signal"
	"syscall"
)

//This method handles OS halt signals.
//The caller can wait for the signal on the returned channel
func OnOSSignal() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGHUP)
	return c
}

//This method handles OS halt signals.
//The given signalCaught function is called once a signal is caught.
//The given onPanic function is called if the handling routine encounters a panic
func OnOSSignalCall(signalCaught func(), onPanic func(r interface{})) {
	c := OnOSSignal()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				onPanic(r)
			}
		}()
		<-c
		signalCaught()
	}()
}
