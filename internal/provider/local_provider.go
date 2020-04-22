package provider

import (
	"errors"
	"fmt"
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

func (p LocalProvider) Get(key string) (io.ReadCloser, error) {
	fileName := path.Join(p.folder, key)
	log.Debug().Msgf("Path: %v", fileName)
	f, err := os.Open(fileName)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrKeyNotFound
	}
	return f, err
}

func (p LocalProvider) Set(key string, in io.Reader) error {
	fileName := path.Join(p.folder, key)
	log.Debug().Msgf("Path: %v", fileName)
	_, err := os.Stat(fileName)
	if err == nil { // already exists, so delete
		os.Remove(fileName)
	}
	f, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error while creating file %v: %w", fileName, err)
	}
	defer f.Close()
	if _, err := io.Copy(f, in); err != nil {
		return fmt.Errorf("error while writing file %v: %w", fileName, err)
	}
	return nil
}