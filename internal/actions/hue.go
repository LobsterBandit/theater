package actions

import (
	"log"

	"github.com/hekmon/plexwebhooks"
)

const (
	HueKind        Kind        = "Philips Hue"
	HueDescription Description = "Control hue lights"
)

// hold some stuff like lights to change, colors to set, etc.
// that can be accessed in execute().
type Hue struct{}

func (h *Hue) kind() Kind {
	return HueKind
}

func (h *Hue) describe() Description {
	return HueDescription
}

func (h *Hue) execute(p interface{}) {
	payload, ok := p.(*plexwebhooks.Payload)
	if !ok {
		return
	}

	log.Printf("Executing %s action in response to event %s\n", h.kind(), payload.Event)
}
