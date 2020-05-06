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
	Folder string `json:"folder"`
}

type LocalProvider struct {
	folder string
}

func NewLocalProvider(conf LocalConf) (LocalProvider, error) {
	p := LocalProvider{folder: conf.Folder}
	if _, err := os.Stat(conf.Folder) ; err != nil {
		return p, fmt.Errorf("error when checking provided root folder: %w", err)
	}
	return p, nil
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
