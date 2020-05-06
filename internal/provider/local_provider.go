package provider

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type LocalConf struct {
	Folder string `json:"folder"`
}

type LocalProvider struct {
	folder string
}

func NewLocalProvider(conf LocalConf) (LocalProvider, error) {
	root, err := filepath.Abs(conf.Folder)
	if err != nil {
		return LocalProvider{}, fmt.Errorf("error while getting absolute path for the root folder: %w", err)
	}
	log.Info().Msgf("Root folder: %v", root)

	p := LocalProvider{folder: root}
	if _, err := os.Stat(root); err != nil {
		return p, fmt.Errorf("error when checking provided root folder: %w", err)
	}

	return p, nil
}

func (p LocalProvider) checkChroot(k string) (string, error) {
	rawP := path.Join(p.folder, k)
	targetedAbs, err := filepath.Abs(rawP)
	if err != nil {
		return "", fmt.Errorf("error while getting absolute path for key %v: %w", k, err)
	}
	if ! strings.HasPrefix(targetedAbs, p.folder) {
		log.Error().Msgf("chroot error, root: %v, targeted: %v", p.folder, targetedAbs)
		return "", ErrChroot
	}
	return targetedAbs, nil

}

func (p LocalProvider) Get(key string) (io.ReadCloser, error) {
	fileName, err := p.checkChroot(key)
	if err != nil {
		return nil, ErrChroot
	}
	log.Debug().Msgf("Path: %v", fileName)

	f, err := os.Open(fileName)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrKeyNotFound
	}

	return f, err
}

func (p LocalProvider) Set(key string, in io.Reader) error {
	fileName, err := p.checkChroot(key)
	if err != nil {
		return ErrChroot
	}
	log.Debug().Msgf("Path: %v", fileName)

	_, err = os.Stat(fileName)
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
