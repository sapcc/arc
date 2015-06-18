package arc

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"
)

type Subprocess struct {
	Command   []string
	cmd       *exec.Cmd
	outPipe   io.ReadCloser
	errPipe   io.ReadCloser
	done      chan struct{}
	outChan   chan string
	exitError error
}

func NewSubprocess(command string, args ...string) *Subprocess {
	return &Subprocess{Command: append([]string{command}, args...)}
}

func ScanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, data[0 : i+1], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
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

func (s *Subprocess) Exited() bool {
	pstate := s.ProcessState().Sys().(syscall.WaitStatus)
	//strangley a signaled proccess on linux is not "Exited()" wtf
	return pstate.Exited() || pstate.Signaled()
}

func (s *Subprocess) Error() error {
	return s.exitError
}

func (s *Subprocess) ProcessState() *os.ProcessState {
	return s.cmd.ProcessState
}

func (s *Subprocess) waitForExit() {
	s.exitError = s.cmd.Wait()
	log.Debugf("Subprocess exited")
	if s.done != nil {
		close(s.done)
	}
}

func (s *Subprocess) scan(pipe io.ReadCloser) {
	scanner := bufio.NewScanner(pipe)
	scanner.Split(ScanLines)
	for scanner.Scan() {
		s.outChan <- scanner.Text()
	}
}
