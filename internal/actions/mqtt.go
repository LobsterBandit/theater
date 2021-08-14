package actions

import (
	"log"

	plex "github.com/hekmon/plexwebhooks"
)

const (
	MqttKind        Kind        = "Philips Hue"
	MqttDescription Description = "Control hue lights"
)

type Mqtt struct {
	topic string
}

func (m *Mqtt) kind() Kind {
	return HueKind
}

func (m *Mqtt) describe() Description {
	return HueDescription
}

func (m *Mqtt) execute(p interface{}) {
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
