// +build linux darwin,!integration

package arc

import (
	"syscall"
	"testing"
	"time"
)

func TestGracefullShutdown(t *testing.T) {

	sub := NewSubprocess("/bin/bash", "-c", `trap 'exit 0' SIGTERM;/bin/sleep 1`)

	lines, err := sub.Start()
	if err != nil {
		t.Error("Failed to start process", err)
		return
	}
	//give the process some time to start
	time.Sleep(50 * time.Millisecond)
	if sub.ProcessState() != nil {
		t.Error("Process already dead", <-lines)
		return
	}

	sub.Kill()
	<-sub.Done()
	pstate := sub.ProcessState().Sys().(syscall.WaitStatus)
	if pstate.ExitStatus() != 0 {
		t.Error("Process didn't exit cleanly")
	}
	if pstate.ExitStatus() != 0 {
		t.Error("Process didn't exit cleanly")
	}

}

func TestForcefullShutdown(t *testing.T) {
	sub := NewSubprocess("/bin/bash", "-c", `trap '' SIGTERM;/bin/sleep 3`)
	//lower the timeout so that the test is not taking longer that neccessarry
	subprocessShutdownTimeout = 100 * time.Millisecond

	lines, err := sub.Start()
	if err != nil {
		t.Error("Failed to start process", err)
		return
	}
	//give the process some time to start
	time.Sleep(50 * time.Millisecond)
	if sub.ProcessState() != nil {
		t.Error("Process already dead", <-lines)
		return
	}

	sub.Kill()
	select {
	case <-sub.Done():
		pstate := sub.ProcessState().Sys().(syscall.WaitStatus)
		if pstate.Signal() != syscall.SIGKILL {
			t.Errorf("Process was %s. Expected %s.", pstate.Signal(), syscall.SIGKILL)
		}
	case <-time.After(1 * time.Second):
		t.Error("Process wasn't killed")
	}

}
