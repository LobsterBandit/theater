package actions

import (
	"log"

	"github.com/hekmon/plexwebhooks"
)

const (
	LoggerKind        Kind        = "Logger"
	LoggerDescription Description = "Log a webhook"
)

type Logger struct {
	log func(p interface{})
}

func (l *Logger) kind() Kind {
	return LoggerKind
}

func (l *Logger) describe() Description {
	return LoggerDescription
}

func (l *Logger) execute(p interface{}) {
	l.log(p)
}

func DefaultLogger() *Logger {
	return &Logger{
		log: func(p interface{}) {
			switch w := p.(type) {
			case *plexwebhooks.Payload:
				log.Printf("received plex webhook event: %s", w.Event)
			default:
				log.Printf("received webhook (%T): %v\n", w, p)
			}
		},
	}
}
