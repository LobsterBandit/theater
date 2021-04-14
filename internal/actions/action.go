package actions

import (
	"log"
)

type ActionKind string

type ActionDescription string

type Action interface {
	kind() ActionKind
	describe() ActionDescription
	execute(p interface{})
}

type ActionHandler struct {
	actions []Action
}

func (w *ActionHandler) Add(a Action) {
	log.Printf("Adding action %s: %s\n", a.kind(), a.describe())
	w.actions = append(w.actions, a)
}

func (w *ActionHandler) ProcessAll(p interface{}) {
	for _, a := range w.actions {
		go a.execute(p)
	}
}

func (w *ActionHandler) ProcessByKind(k ActionKind, p interface{}) {
	for _, a := range w.actions {
		if a.kind() == k {
			go a.execute(p)
		}
	}
}
