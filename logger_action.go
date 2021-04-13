package main

import (
	"log"

	"github.com/hekmon/plexwebhooks"
)

const (
	LoggerActionKind        ActionKind        = "Logger"
	LoggerActionDescription ActionDescription = "Log a webhook"
)

type LoggerAction struct {
	log func(p interface{})
}

func (l *LoggerAction) kind() ActionKind {
	return LoggerActionKind
}

func (l *LoggerAction) describe() ActionDescription {
	return LoggerActionDescription
}

func (l *LoggerAction) execute(p interface{}) {
	l.log(p)
}

func DefaultLogAction() *LoggerAction {
	return &LoggerAction{
		log: func(p interface{}) {
			switch w := p.(type) {
			case *plexwebhooks.Payload:
				log.Printf("received plex webhook event: %s", p)
			default:
				log.Printf("received webhook (%T): %v\n", w, p)
			}
		},
	}
}
