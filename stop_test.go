package stop

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
)

func TestSendStopSignalAndWait(t *testing.T) {
	tempDir := os.TempDir()
	os.Remove(path.Join(tempDir, stopProcessFileName))
	defer os.Remove(path.Join(tempDir, stopProcessFileName))
	processStoppedPath := path.Join(tempDir, processStoppedFileName)
	os.Remove(processStoppedPath)
	defer os.Remove(processStoppedPath)

	stopChannel := make(chan string)
	SendStopSignalAndWait(stopChannel)

	ioutil.WriteFile(processStoppedPath, []byte{}, 600)

	select {
	case stopMessage := <-stopChannel:
		if stopMessage != Signal {
			t.Error("Unexpected stop signal. Received", stopMessage, "but expected", Signal)
		}
	case <-time.After(5 * time.Second):
		t.Error("Timed out waiting for process to stop")
	}
}

func TestSendStopSignalAndTimeout(t *testing.T) {
	tempDir := os.TempDir()
	os.Remove(path.Join(tempDir, stopProcessFileName))
	defer os.Remove(path.Join(tempDir, stopProcessFileName))
	processStoppedPath := path.Join(tempDir, processStoppedFileName)
	os.Remove(processStoppedPath)
	defer os.Remove(processStoppedPath)

	stopChannel := make(chan string)
	SendStopSignalAndWait(stopChannel)

	select {
	case stopMessage := <-stopChannel:
		t.Error("Process should not have stopped but received signal", stopMessage)
	case <-time.After(2 * time.Second):

	}
}

func TestReceiveStopSignal(t *testing.T) {
	tempDir := os.TempDir()
	stopProcessPath := path.Join(tempDir, stopProcessFileName)
	os.Remove(stopProcessPath)
	defer os.Remove(stopProcessPath)

	stopChannel := ListenForStopSignal()

	ioutil.WriteFile(stopProcessPath, []byte{}, 600)

	select {
	case stopMessage := <-stopChannel:
		if stopMessage != Signal {
			t.Error("Unexpected stop signal. Received", stopMessage, "but expected", Signal)
		}
	case <-time.After(5 * time.Second):
		t.Error("Timed out waiting for process to stop")
	}
}

func TestProcessStoppedSignalIsCreated(t *testing.T) {
	tempDir := os.TempDir()
	processStoppedPath := path.Join(tempDir, processStoppedFileName)
	os.Remove(processStoppedPath)
	defer os.Remove(processStoppedPath)

	SignalThatProcessHasStopped()
	_, err := os.Stat(processStoppedPath)
	if err != nil {
		t.Errorf("Process stopped signal file should have been created at %s but received error %v", processStoppedPath, err)
	}
}
