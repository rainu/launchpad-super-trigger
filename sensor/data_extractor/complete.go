package data_extractor

type Complete struct {
}

func (c Complete) Extract(data []byte) ([]byte, error) {
	return data, nil
}
