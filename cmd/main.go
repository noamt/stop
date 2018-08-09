package main

import (
	"log"
	"time"
	"github.com/noamt/stop"
)

func main() {
	stopChannel := make(chan string)
	stop.SendStopSignalAndWait(stopChannel)
	select {
	case stopMessage := <-stopChannel:
		log.Println(stopMessage)
	case <-time.After(30 * time.Second):
		log.Fatal("Timed out waiting for process to stop")
	}
}
