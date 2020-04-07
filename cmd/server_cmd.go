package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/barasher/file-server/internal/provider"
	"github.com/barasher/file-server/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
)

var ServerCmd = &cobra.Command{
	Use:   "s3pub",
	Short: "S3 bucket publisher",
	RunE:  execute,
	SilenceUsage:true,
}
var confFile string

const (
	s3ProviderID = "S3"
	localProviderID="local"
)

func init() {
	ServerCmd.Flags().StringVarP(&confFile, "file", "f", "", "Configuration file")
	ServerCmd.MarkFlagRequired("file")
}

func loadConfFile(path string) (ServerConf, error) {
	conf := ServerConf{}
	confReader, err := os.Open(confFile)
	if err != nil {
		return conf, fmt.Errorf("Error while opening configuration file %v :%v", confFile, err)
	}
	defer confReader.Close()
	err = json.NewDecoder(confReader).Decode(&conf)
	if err != nil {
		return conf, fmt.Errorf("Error while unmarshaling configuration file %v :%v", confFile, err)
	}
	return conf, nil
}

func setLoggingLevel(lvl string) error {
	if lvl != "" {
		lvl, err := zerolog.ParseLevel(lvl)
		if err != nil {
			return fmt.Errorf("error while setting logging level (%v): %w", lvl, err)
		}
		zerolog.SetGlobalLevel(lvl)
		log.Debug().Msgf("Logging level: %v", lvl)
	}
	return nil
}

func execute(cmd *cobra.Command, args []string) error {
	var err error

	var conf ServerConf
	if conf, err = loadConfFile(confFile); err != nil {
		return err
	}

	if err := setLoggingLevel(conf.LoggingLevel); err != nil {
		return fmt.Errorf("error while setting logging level (%v): %w", conf.LoggingLevel, err)
	}

	var prov provider.Provider
	switch conf.Type {
	case s3ProviderID:
		log.Info().Msg("Provider: S3")
		if prov, err = provider.NewS3Provider(conf.S3Conf); err != nil {
			return err
		}
	case localProviderID:
		log.Info().Msg("Provider: local")
		if prov, err = provider.NewLocalProvider(conf.LocalConf); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unknown provider type (%v)", conf.Type)
	}

	server.Run(prov)
	return nil
}


type ServerConf struct { // TODO normaliser les balises
	Type         string
	S3Conf       provider.S3Conf    `json:"S3"`
	LocalConf    provider.LocalConf `json:"Local"`
	LoggingLevel string
	Port uint
}

