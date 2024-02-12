package kafka

type MessageProcessor interface {
	ProcessMessage(message []byte) error
	ProcessMessages(messages [][]byte) error
}
