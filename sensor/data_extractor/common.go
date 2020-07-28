package data_extractor

type Extractor interface {
	Extract(data []byte) ([]byte, error)
}
