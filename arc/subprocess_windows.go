package arc

func (s *Subprocess) Kill() {
	s.cmd.Process.Kill()
}
