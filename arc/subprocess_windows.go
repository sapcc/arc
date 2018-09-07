package arc

import (
	"fmt"
	"os/exec"
)

func (s *Subprocess) Kill() {
	s.cmd.Process.Kill()
}

func (s *Subprocess) prepareCmd() *exec.Cmd {
	cmd := exec.Command(s.Command[0], s.Command[1:]...) // #nosec
	fmt.Printf("%+v\n", cmd)
	return cmd
}
