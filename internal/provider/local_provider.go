package provider

import (
	"errors"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path"
)

type LocalConf struct {
	Folder     string
}

type LocalProvider struct {
	folder	string
}

func NewLocalProvider(conf LocalConf) (LocalProvider, error) {
	return LocalProvider{folder:conf.Folder}, nil
}

func (p LocalProvider) Provide(key string) (io.ReadCloser, error) {
	fileName := path.Join(p.folder, key)
	log.Info().Msgf("Path: %v", fileName)
	f, err := os.Open(fileName)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrKeyNotFound
	}
	return f, err
}

func (p LocalProvider) Close() {
	// Nothing
}