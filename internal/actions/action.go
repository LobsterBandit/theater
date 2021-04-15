package actions

import (
	"log"
)

type Kind string

type Description string

type Action interface {
	kind() Kind
	describe() Description
	execute(p interface{})
}

type Handler struct {
	actions []Action
}

func (w *Handler) Add(a ...Action) {
	for _, action := range a {
		log.Printf("Adding action %s: %s\n", action.kind(), action.describe())
		w.actions = append(w.actions, action)
	}
}

func (w *Handler) ProcessAll(p interface{}) {
	for _, a := range w.actions {
		go a.execute(p)
	}
}

func (w *Handler) ProcessByKind(k Kind, p interface{}) {
	for _, a := range w.actions {
		if a.kind() == k {
			go a.execute(p)
		}
	}
}
