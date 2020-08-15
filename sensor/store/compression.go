package store

import (
	"bytes"
	"compress/gzip"
	"io"
)

type GzipCompressed struct {
	Delegate Store
}

func (g *GzipCompressed) Set(data []byte) error {
	if len(data) == 0 {
		return g.Delegate.Set(data)
	}

	var compressed bytes.Buffer
	writer := gzip.NewWriter(&compressed)

	//compress the data
	if _, err := writer.Write(data); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	//save the compressed data
	return g.Delegate.Set(compressed.Bytes())
}

func (g *GzipCompressed) Get() ([]byte, error) {
	compressed, err := g.Delegate.Get()
	if err != nil {
		return compressed, err
	}
	if len(compressed) == 0 {
		return compressed, nil
	}

	reader, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, err
	}

	//decompress the data
	var decompressed bytes.Buffer

	if _, err := io.Copy(&decompressed, reader); err != nil {
		return nil, err
	}
	if err := reader.Close(); err != nil {
		return nil, err
	}

	//return the decompressed data
	return decompressed.Bytes(), nil
}
