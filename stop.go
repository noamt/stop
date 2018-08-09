package stop

import (
	"log"
	"os"
	"path/filepath"
	"io/ioutil"
	"time"
)

const stopProcessFileName = "process.stop"
const processStoppedFileName = "process.stopped"
const Signal = "stop"

func SendStopSignalAndWait(stopChannel chan string) {
	sendStopSignal()
	go waitForProcessToStop(stopChannel)
}

func ListenForStopSignal(c chan<- string) {
	go func() {
		log.Println("Starting stop signal listener")
		checkForSignal := true
		tempDir := os.TempDir()
		stopProcessFile := filepath.Join(tempDir, stopProcessFileName)

		for checkForSignal {
			if _, err := os.Stat(stopProcessFile); err == nil {
				log.Println("Found stop signal:", stopProcessFile)
				c <- Signal
				checkForSignal = false
			}
		}
	}()
}

func SignalThatProcessHasStopped() {
	tempDir := os.TempDir()
	processStoppedFile := filepath.Join(tempDir, processStoppedFileName)
	ioutil.WriteFile(processStoppedFile, []byte{}, 600)
}

func sendStopSignal() {
	tempDir := os.TempDir()
	stopProcessFile := filepath.Join(tempDir, stopProcessFileName)
	log.Println("Received stop command. Writing message to", stopProcessFile)
	ioutil.WriteFile(stopProcessFile, []byte{}, 600)
}

func waitForProcessToStop(stopChannel chan string) {
	log.Println("Waiting for process to stop")
	tempDir := os.TempDir()
	processStoppedFile := filepath.Join(tempDir, processStoppedFileName)
	var stopped bool
	for !stopped {
		time.Sleep(time.Second * 2)
		_, err := os.Stat(processStoppedFile)
		stopped = err == nil
	}
	stopChannel <- Signal
}
