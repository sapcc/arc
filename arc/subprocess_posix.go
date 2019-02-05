// +build linux darwin

package arc

import (
	"os/exec"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
)

var subprocessShutdownTimeout = 2 * time.Second

func (s *Subprocess) Kill() {
	doneChan := s.Done()
	if err := s.cmd.Process.Signal(syscall.SIGTERM); err != nil {
		log.Errorf("Error killing process: %s\n", err)
	}
	select {
	case <-doneChan:
	case <-time.After(subprocessShutdownTimeout):
		log.Warnf("Process didn't terminate gracefully within %v. Killing it.", subprocessShutdownTimeout)
		//We kill the entire process group here to make sure child processes of the process die as well
		if err := syscall.Kill(-s.cmd.Process.Pid, syscall.SIGKILL); err != nil {
			log.Errorf("Error killing the entire process: %s\n", err)
		}
	}
}

func (s *Subprocess) prepareCmd() *exec.Cmd {
	cmd := exec.Command(s.Command[0], s.Command[1:]...) // #nosec
	if s.Env != nil {
		cmd.Env = s.Env
	}
	cmd.Dir = s.Dir
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	return cmd
}
