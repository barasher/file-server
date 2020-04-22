package internal

import (
	"encoding/json"
	"fmt"
	"github.com/barasher/file-server/internal/provider"
	"os"
)

const (
	S3ProviderID    = "s3"
	LocalProviderID = "local"
)

type ServerConf struct {
	Type         string
	S3Conf       provider.S3Conf    `json:"S3"`
	LocalConf    provider.LocalConf `json:"Local"`
	LoggingLevel string
	Port uint
}

func LoadConfFile(path string) (ServerConf, error) {
	conf := ServerConf{}
	confReader, err := os.Open(path)
	if err != nil {
		return conf, fmt.Errorf("Error while opening configuration file %v :%v", path, err)
	}
	defer confReader.Close()
	err = json.NewDecoder(confReader).Decode(&conf)
	if err != nil {
		return conf, fmt.Errorf("Error while unmarshaling configuration file %v :%v", path, err)
	}
	return conf, nil
}