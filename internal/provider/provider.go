package provider

import (
	"fmt"
	"io"
)

type Provider interface {
	Provide(key string) (io.ReadCloser, error)
	Close()
}

var ErrKeyNotFound error = fmt.Errorf("Key not found")
