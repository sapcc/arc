package network

type Source struct{}

func New() Source {
	return Source{}
}

func newFacts() map[string]interface{} {
	facts := make(map[string]interface{})
	facts["ipaddress"] = nil
	facts["macaddress"] = nil
	facts["default_interface"] = nil
	facts["default_gateway"] = nil

	return facts
}

func (h Source) Name() string {
	return "network"
}
