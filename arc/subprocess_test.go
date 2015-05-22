package arc

import (
	"syscall"
	"testing"
	"time"
)

func TestCommandNotFound(t *testing.T) {

	sub := NewSubprocess("i-dont-exist")

	lines, err := sub.Start()
	if lines != nil {
		t.Error("returned channel should be nil")
	}
	if err == nil {
		t.Error("error should be non nil")
	}

}

func TestSubprocess(t *testing.T) {
	sub := NewSubprocess("echo", "tut")

	lines, _ := sub.Start()

	output := <-lines
	if output != "tut" {
		t.Error("Unexpected output: ", output)
	}
	<-sub.Done()
}

func TestKillSubprocess(t *testing.T) {

	sub := NewSubprocess("sleep", "2")

	_, err := sub.Start()
	if err != nil {
		t.Errorf("Error starting process", err)
	}

	//five the process some time to start
	time.Sleep(50 * time.Millisecond)
	sub.Kill()
	select {
	case <-sub.Done():
		t.Errorf(sub.ProcessState().Sys().(syscall.WaitStatus))
		if !sub.ProcessState().Exited() {
			t.Error("Process did not exit", sub.ProcessState().Sys().(syscall.WaitStatus).Signaled())
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("process didn't terminate")
	}

}
