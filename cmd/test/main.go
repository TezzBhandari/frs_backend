package main

import (
	"flag"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	fileInputPath  string
	fileOutputPath string
	debug          bool
)

func init() {
	flag.StringVar(&fileInputPath, "input", "", "Path to the input file")
	flag.StringVar(&fileOutputPath, "output", "", "Path to the output file")
	flag.BoolVar(&debug, "debug", false, "Set log level to debug")

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func main() {
	flag.Parse()
	if fileInputPath == "" {
		log.Error().Msg("input file path is required")
		flag.Usage()
		os.Exit(1)
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Debug().Msg("This message appear only when log level set to Debug")
	log.Info().Msg(("This message appear when log level set to Debug or Info"))
}
