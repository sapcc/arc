package arc

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type Subprocess struct {
	Command []string
	cmd     *exec.Cmd
	outPipe io.ReadCloser
	errPipe io.ReadCloser
	done    chan struct{}
	outChan chan string
}

func NewSubprocess(command string, args ...string) *Subprocess {
	return &Subprocess{Command: append([]string{command}, args...)}
}

func (s *Subprocess) Start() (<-chan string, error) {

	s.cmd = exec.Command(s.Command[0], s.Command[1:]...)
	outPipe, err := s.cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	s.outPipe = outPipe
	errPipe, err := s.cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	s.errPipe = errPipe

	if err := s.cmd.Start(); err != nil {
		return nil, err
	}
	log.Debugf("Started subprocess %s", strings.Join(s.Command, " "))
	s.outChan = make(chan string, 10)

	go s.scan(s.outPipe)
	go s.scan(s.errPipe)
	go s.waitForExit()

	return s.outChan, nil

}

func (s *Subprocess) Done() <-chan struct{} {
	if s.done == nil {
		s.done = make(chan struct{})
	}
	return s.done
}

func (s *Subprocess) ProcessState() *os.ProcessState {
	return s.cmd.ProcessState
}

func (s *Subprocess) waitForExit() {
	s.cmd.Wait()
	log.Debugf("Subprocess exited")
	if s.done != nil {
		close(s.done)
	}
}

func (s *Subprocess) scan(pipe io.ReadCloser) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		s.outChan <- scanner.Text()
	}
}
