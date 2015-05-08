package onos

type Request struct {
	Version   int
	Sender    string
	RequestID string
	To        string
	Type      string
	Timeout   uint64
	Agent     string
	Action    string
	Payload   string
}

type Reply struct {
	Version   int
	Sender    string
	RequestID string
	Type      string
	Agent     string
	Action    string
	State     string
	Final     bool
	Payload   string
}

func CreateReply(request *Request, payload string) *Reply {

	return &Reply{
		Version: 1,
		Type:    "reply",
		Agent:   request.Agent,
		Action:  request.Action,
		Payload: payload,
	}

}
