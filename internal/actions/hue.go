package actions

import (
	"log"

	"github.com/amimof/huego"
	plex "github.com/hekmon/plexwebhooks"
)

const (
	HueKind        Kind        = "Philips Hue"
	HueDescription Description = "Control hue lights"
)

type Hue struct {
	Bridge     *huego.Bridge
	Group      int
	Scene      string
	PlexEvents map[plex.EventType]struct{}
	PlexPlayer string
	PlexUser   string
}

func (h *Hue) kind() Kind {
	return HueKind
}

func (h *Hue) describe() Description {
	return HueDescription
}

func (h *Hue) execute(p interface{}) {
	payload, ok := p.(*plex.Payload)
	if !ok {
		return
	}

	// Only act on events matching the desired user and player
	if payload.Account.Title != h.PlexUser || payload.Player.Title != h.PlexPlayer {
		return
	}

	// Only act on events matching the desired event types
	if _, ok := h.PlexEvents[payload.Event]; !ok {
		return
	}

	log.Printf("Executing %s action in response to event %s\n", h.kind(), payload.Event)

	if _, err := h.Bridge.RecallScene(h.Scene, h.Group); err != nil {
		log.Println(err)
	}
}
