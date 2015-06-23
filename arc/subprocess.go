package arc

import (
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

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
	wg        sync.WaitGroup
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
	s.wg.Wait()
	if s.done != nil {
		close(s.done)
	}
}

func (s *Subprocess) scan(pipe io.ReadCloser) {
	s.wg.Add(1)
	defer s.wg.Done()

	chunker := NewChunkedReader(pipe, 1*time.Second, 4096)

	for {
		chunk, err := chunker.Read()
		if chunk != nil {
			log.Debugf("Sending chunk (size: %d)", len(chunk))
			s.outChan <- string(chunk)
		}
		if err != nil {
			return
		}
	}
}
