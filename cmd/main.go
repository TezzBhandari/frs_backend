package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/TezzBhandari/frs/http"
	"github.com/TezzBhandari/frs/postgres"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	addr  string
	debug bool
	dsn   string
)

func init() {
	flag.StringVar(&addr, "addr", "", "Specifies the tcp server address for server to listen on")
	flag.BoolVar(&debug, "debug", false, "Sets log level flag to default")
	flag.StringVar(&dsn, "dsn", "", "Sets database dsn")

	flag.Parse()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("Log level set to Debug")
	}

	if dsn == "" {
		log.Info().Msg("Set -dsn flag")
		os.Exit(1)
	}

	if addr == "" {
		log.Info().Msg("Set -addr flag")
		os.Exit(1)
	}

}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt)
	go func() {
		<-sigChan
		cancel()
	}()

	m := NewMain()
	if err := m.run(); err != nil {
		log.Error().Err(err).Msg("")
		if err = m.close(); err != nil {
			log.Error().Err(err).Msg("")
		}
		os.Exit(1)
	}

	<-ctx.Done()

	log.Info().Msg("interrupt triggered")

	if err := m.close(); err != nil {
		log.Error().Err(err).Msg("")
		os.Exit(1)
	}

}

type Main struct {
	HttpServer *http.Server
	DB         *postgres.DB
}

func NewMain() *Main {
	return &Main{
		HttpServer: http.NewHttpServer(),
		DB:         postgres.NewDB(dsn),
	}
}

func (m *Main) run() error {
	m.HttpServer.Addr = addr

	if err := m.DB.Open(); err != nil {
		return fmt.Errorf("cannot open db: %w", err)
	}

	userService := postgres.NewUserService(m.DB)
	fundRaiserService := postgres.NewFundRaiserService(m.DB)

	// attach underlying services to http server
	m.HttpServer.UserService = userService
	m.HttpServer.FundRaiserService = fundRaiserService

	if err := m.HttpServer.Open(); err != nil {
		return fmt.Errorf("cannot start server: %w", err)
	}

	fmt.Printf("running: url=%q dsn=%q\n", m.HttpServer.Url(), dsn)

	return nil

}

func (m *Main) close() error {
	if err := m.DB.Close(); err != nil {
		return fmt.Errorf("closing db error: %w", err)
	}
	if err := m.HttpServer.Close(); err != nil {
		return fmt.Errorf("error closinng server: %w", err)
	}
	return nil
}
