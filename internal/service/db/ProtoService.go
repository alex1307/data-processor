package service

type ProtoService interface {
	Save(message []byte) (uint64, error)
	SaveAll(messages [][]byte) error
}
