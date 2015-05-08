package onos

type Message struct {
	Version   int
	Sender    string
	RequestID string
	Type      string
	Timeout   uint64
	Agent     string
	Action    string
	Payload   string
}

func CreateReply(request *Message, payload string) *Message {

	return &reply{
		Version: 1,
		Type:    "reply",
		Agent:   request.Agent,
		Action:  request.Action,
		Payload: payload,
	}

}
