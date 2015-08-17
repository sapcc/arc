package service

var serviceName = "arc"
var serviceDisplayName = "Arc Agent" //mostly for windows
var serviceDescription = "Monsoon remote control agent"

type Service interface {
	Install() error
	Status() (string, error)
	Restart() error
	Start() error
	Stop() error
}

type service struct {
	dir string
}

func New(dir string) Service {
	return service{dir: dir}
}
