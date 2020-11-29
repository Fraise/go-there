package logging

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go-there/config"
	"go-there/data"
	"io"
	"os"
)

// Init initialize the logger from the provided configuration and returns a non-nil file if one was open or created.
// Returns data.ErrInit if the initialization fails.
func Init(conf *config.Configuration) (*os.File, error) {
	var output io.Writer
	var f *os.File

	switch conf.Logs.File {
	case "", "$stdout":
		output = os.Stdout
	case "$stderr":
		output = os.Stderr
	default:
		var err error
		f, err = os.OpenFile(conf.Logs.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			return nil, fmt.Errorf("%w : %s", data.ErrInit, err)
		}

		output = f
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if conf.Logs.AsJSON {
		log.Logger = log.Output(output)
	} else {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: output})
	}

	return f, nil
}
