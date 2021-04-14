package actions

import (
	"log"

	"github.com/amimof/huego"
	"github.com/hekmon/plexwebhooks"
)

const (
	HueKind        Kind        = "Philips Hue"
	HueDescription Description = "Control hue lights"
)

type Hue struct {
	Bridge     *huego.Bridge
	PlexEvent  plexwebhooks.EventType
	PlexPlayer string
	PlexUser   string
	Lights     map[int]huego.State
}

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

	// Only act on events matching the desired user and player
	if h.PlexUser != payload.Account.Title ||
		h.PlexPlayer != payload.Player.Title ||
		h.PlexEvent != payload.Event {
		return
	}

	log.Printf("Executing %s action in response to event %s\n", h.kind(), payload.Event)

	for i, l := range h.Lights {
		resp, err := h.Bridge.SetLightState(i, l)
		if err != nil {
			log.Println(err)
		}

		log.Printf("%v\n", resp)
	}
}
