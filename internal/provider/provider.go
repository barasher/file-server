package provider

import (
	"fmt"
	"io"
)

type Provider interface {
	Get(key string) (io.ReadCloser, error)
	Set(key string, b io.Reader) error
}

var ErrKeyNotFound error = fmt.Errorf("Key not found")
var ErrChroot error = fmt.Errorf("Chroot error")
