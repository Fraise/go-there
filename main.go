package main

import (
	"context"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go-there/api"
	"go-there/config"
	"go-there/datasource"
	"go-there/gopath"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var configPath = flag.String("config", "", "Path to the configuration file")
	flag.Parse()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	conf, err := config.Init(*configPath)

	if err != nil {
		log.Fatal().Err(err).Send()
	}

	e := gin.New()

	e.Use(gin.Logger())
	e.Use(gin.Recovery())

	ds, err := datasource.Init(conf)

	if err != nil {
		log.Fatal().Err(err).Send()
	}

	gopath.Init(e, ds)
	api.Init(e, ds)

	s := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: e,
	}

	go func() {
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			log.Error().Err(err).Send()
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	sig := <-c

	log.Info().Msgf("received %v", sig)
	log.Info().Msgf("shutting down the http server")

	if err := s.Shutdown(context.Background()); err != nil {
		log.Error().Err(err).Msg("error shutting down the http server")
	}
}
