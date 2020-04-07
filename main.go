package main

import (
	"github.com/barasher/file-server/cmd"
	"github.com/rs/zerolog/log"
	"os"
)

const (
	retOk int = 0
	retKo int = 1
)

func main() {
	if err := cmd.ServerCmd.Execute(); err != nil {
		log.Error().Msgf("%v", err)
		os.Exit(retKo)
	}
	os.Exit(retOk)
}
