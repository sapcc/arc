// +build linux darwin

package arc

import (
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
)

var subprocessShutdownTimeout = 2 * time.Second

func (s *Subprocess) Kill() {
	doneChan := s.Done()
	s.cmd.Process.Signal(syscall.SIGTERM)
	select {
	case <-doneChan:
	case <-time.After(subprocessShutdownTimeout):
		log.Warn("Process didn't terminate gracefully within 2 seconds. Killing it.")
		s.cmd.Process.Kill()
	}
}
