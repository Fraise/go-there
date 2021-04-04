package main

import (
	"context"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go-there/api"
	"go-there/cache"
	"go-there/config"
	"go-there/database"
	"go-there/datasource"
	"go-there/gopath"
	"go-there/health"
	"go-there/logging"
	"go-there/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var configPath = flag.String("config", "", "Path to the configuration file")
	flag.Parse()

	// Basic logging for the initialization
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	conf, err := config.Init(*configPath)

	if err != nil {
		log.Fatal().Err(err).Send()
	}

	err = api.ApplyUserSettings(conf)

	if err != nil {
		log.Fatal().Err(err).Send()
	}

	logFile, err := logging.Init(conf)

	if err != nil {
		log.Fatal().Err(err).Send()
	}

	if logFile != nil {
		defer func() {
			_ = logFile.Close()
		}()
	}

	if conf.Server.Mode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	e := gin.New()

	e.Use(gin.Logger())
	e.Use(gin.Recovery())

	db, err := database.Init(conf)

	if err != nil {
		log.Fatal().Err(err).Send()
	}

	ds := datasource.Init(db, cache.Init(conf))

	health.Init(conf, e)
	gopath.Init(conf, e, ds)
	api.Init(conf, e, ds)

	if conf.Server.HttpListenPort <= 0 && conf.Server.HttpsListenPort <= 0 {
		log.Fatal().Err(err).Msg("no listening port configured")
	}

	s, tlsServer := server.Start(conf, e)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	sig := <-c

	log.Info().Msgf("received %v", sig)
	log.Info().Msgf("shutting down the http server")

	if s != nil {
		if err := s.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msg("error shutting down the http server")
		}
	}

	if tlsServer != nil {
		if err := tlsServer.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msg("error shutting down the http server")
		}
	}

}
