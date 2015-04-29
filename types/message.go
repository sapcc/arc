package types

type Message struct {
	Version   int
	RequestId string
	Type      string
	Agent     string
	Action    string
	Payload   string
}
