package stop

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

const stopProcessFileName = "process.stop"
const processStoppedFileName = "process.stopped"
const Signal = "stop"

//This method is called by the custom stop command.
//It expects to receive a channel on which the caller can wait.
//It return any error that might occur while sending the signal
func SendStopSignalAndWait(stopChannel chan string) error {
	signalError := sendStopSignal()
	if signalError != nil {
		return signalError
	}
	go waitForProcessToStop(stopChannel)
	return nil
}

//This method should be called by the application that wants to handle the halt process.
//It returns a channel on which the caller can wait, until a stop signal is found.
func ListenForStopSignal() <-chan string {
	c := make(chan string)
	go func() {
		log.Println("Starting stop signal listener")
		checkForSignal := true
		tempDir := os.TempDir()
		stopProcessFile := filepath.Join(tempDir, stopProcessFileName)

		for checkForSignal {
			time.Sleep(time.Second * 2)
			if _, err := os.Stat(stopProcessFile); err == nil {
				log.Println("Found stop signal:", stopProcessFile)
				c <- Signal
				checkForSignal = false
			} else if !os.IsNotExist(err) {
				log.Fatalf("Error while testing for the existence of the stop signal %s: %v", stopProcessFile, err)
			}
		}
	}()
	return c
}

//This method should be called by the application once it has finished cleaning up, and is ready to shutdown
func SignalThatProcessHasStopped() error {
	tempDir := os.TempDir()
	processStoppedFile := filepath.Join(tempDir, processStoppedFileName)
	return ioutil.WriteFile(processStoppedFile, []byte{}, 600)
}

func sendStopSignal() error {
	tempDir := os.TempDir()
	stopProcessFile := filepath.Join(tempDir, stopProcessFileName)
	log.Println("Received stop command. Writing message to", stopProcessFile)
	return ioutil.WriteFile(stopProcessFile, []byte{}, 600)
}

func waitForProcessToStop(stopChannel chan string) {
	log.Println("Waiting for process to stop")
	tempDir := os.TempDir()
	processStoppedFile := filepath.Join(tempDir, processStoppedFileName)
	var stopped bool
	for !stopped {
		time.Sleep(time.Second * 2)
		if _, err := os.Stat(processStoppedFile); err == nil {
			stopped = true
		} else if !os.IsNotExist(err) {
			log.Fatalf("Error while testing for the existence of the stopped signal %s: %v", processStoppedFile, err)
		}
	}
	stopChannel <- Signal
}
