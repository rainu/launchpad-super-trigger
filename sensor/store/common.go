package store

type Store interface {
	Set(data []byte) error
	Get() ([]byte, error)
}
