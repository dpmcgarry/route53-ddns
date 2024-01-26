package main

import (
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	*zerolog.Logger
}

func configureLogging() *Logger {
	var writers []io.Writer

	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}

	writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	writers = append(writers, newRollingFile())
	mw := io.MultiWriter(writers...)

	// zerolog.SetGlobalLevel(zerolog.DebugLevel)
	logger := zerolog.New(mw).With().Timestamp().Logger()
	// Adds the ability to change the logging level using an environment variable
	level := strings.ToUpper(os.Getenv("LOGLEVEL"))
	switch level {
	case "TRACE":
		log.Logger = zerolog.New(mw).
			Level(zerolog.TraceLevel).
			With().
			Timestamp().
			Caller().
			Logger()
	case "DEBUG":
		log.Logger = zerolog.New(mw).
			Level(zerolog.DebugLevel).
			With().
			Timestamp().
			Caller().
			Logger()
	case "INFO":
		log.Logger = zerolog.New(mw).
			Level(zerolog.InfoLevel).
			With().
			Timestamp().
			Caller().
			Logger()
	case "WARN":
		log.Logger = zerolog.New(mw).
			Level(zerolog.WarnLevel).
			With().
			Timestamp().
			Caller().
			Logger()
	case "ERROR":
		log.Logger = zerolog.New(mw).
			Level(zerolog.ErrorLevel).
			With().
			Timestamp().
			Caller().
			Logger()
	default: // Default to info level
		log.Logger = zerolog.New(mw).
			Level(zerolog.DebugLevel).
			With().
			Timestamp().
			Caller().
			Logger()
	}

	return &Logger{
		Logger: &logger,
	}
}

func newRollingFile() io.Writer {
	return &lumberjack.Logger{
		Filename:   path.Join(viper.GetString("logdir") + "route53-ddns.log"),
		MaxBackups: 5,
		MaxSize:    20,
		MaxAge:     30,
	}
}
