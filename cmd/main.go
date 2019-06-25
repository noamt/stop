package main

import (
	"github.com/noamt/stop"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	stopChannel := make(chan string)
	signalError := stop.SendStopSignalAndWait(stopChannel)
	if signalError != nil {
		log.Fatalf("Error while signalling to stop %v", signalError)
	}
	select {
	case stopMessage := <-stopChannel:
		log.Println(stopMessage)
	case <-time.After(secondsToWait() * time.Second):
		log.Fatal("Timed out waiting for process to stop")
	}
}

func secondsToWait() time.Duration {
	timeout := os.Getenv("STOP_TIMEOUT_SECONDS")
	if timeout != "" {
		timeoutInt, e := strconv.Atoi(timeout)
		if e == nil {
			return time.Duration(timeoutInt)
		}
	}
	return 30
}
