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
	Type         string `json:"type"`
	S3Conf       provider.S3Conf    `json:"s3"`
	LocalConf    provider.LocalConf `json:"local"`
	LoggingLevel string `json:"loggingLevel"`
	Port uint `json:"port"`
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