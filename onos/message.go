package onos

type Message struct {
	Version   int
	Sender    string
	RequestID string
	Type      string
	Agent     string
	Action    string
	Payload   string
}
