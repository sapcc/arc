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
