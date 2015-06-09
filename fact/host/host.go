package host

type Source struct{}

func New() Source {
	return Source{}
}

func (h Source) Name() string {
	return "host"
}
