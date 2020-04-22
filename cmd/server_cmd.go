package cmd

import (
	"fmt"
	"github.com/barasher/file-server/internal"
	"github.com/barasher/file-server/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var ServerCmd = &cobra.Command{
	Use:          "s3pub",
	Short:        "S3 bucket publisher",
	RunE:         execute,
	SilenceUsage: true,
}
var confFile string



func init() {
	ServerCmd.Flags().StringVarP(&confFile, "file", "f", "", "Configuration file")
	ServerCmd.MarkFlagRequired("file")
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

	var conf internal.ServerConf
	if conf, err = internal.LoadConfFile(confFile); err != nil {
		return err
	}

	if err := setLoggingLevel(conf.LoggingLevel); err != nil {
		return fmt.Errorf("error while setting logging level (%v): %w", conf.LoggingLevel, err)
	}

	s, err := server.NewServer(conf)
	if err != nil {
		return fmt.Errorf("error while initializing server: %w", err)
	}

	s.Run()

	return nil
}
