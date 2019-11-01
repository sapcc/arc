package service

func (s service) Install(string) error {
	panic("Not implemented on this platform")
}

func (s service) Status() (State, string, error) {
	panic("Not implemented on this platform")
}
func (s service) Restart() error {
	panic("Not implemented on this platform")
}
func (s service) Start() error {
	panic("Not implemented on this platform")
}
func (s service) Stop() error {
	panic("Not implemented on this platform")
}
