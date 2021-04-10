package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/hekmon/plexwebhooks"
)

var webhookBaseURL = "https://discordapp.com/api/webhooks/" //nolint:gochecknoglobals

type Webhook struct {
	ID      string
	Token   string
	Message *Message
}

// A Message stores all data related to a specific Discord message.
type Message struct {
	// A list of thumbnails received from a plex webhook.
	Images []*plexwebhooks.Thumbnail
	// A list of embeds present in the message.
	Embeds []*MessageEmbed `json:"embeds"`
}

// MessageEmbedFooter is a part of a MessageEmbed struct.
type MessageEmbedFooter struct {
	Text         string `json:"text,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// MessageEmbedImage is a part of a MessageEmbed struct.
type MessageEmbedImage struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

// MessageEmbedThumbnail is a part of a MessageEmbed struct.
type MessageEmbedThumbnail struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

// MessageEmbedVideo is a part of a MessageEmbed struct.
type MessageEmbedVideo struct {
	URL    string `json:"url,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

// MessageEmbedProvider is a part of a MessageEmbed struct.
type MessageEmbedProvider struct {
	URL  string `json:"url,omitempty"`
	Name string `json:"name,omitempty"`
}

// MessageEmbedAuthor is a part of a MessageEmbed struct.
type MessageEmbedAuthor struct {
	URL          string `json:"url,omitempty"`
	Name         string `json:"name,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// MessageEmbedField is a part of a MessageEmbed struct.
type MessageEmbedField struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

// An MessageEmbed stores data for message embeds.
type MessageEmbed struct {
	URL         string                 `json:"url,omitempty"`
	Type        string                 `json:"type,omitempty"`
	Title       string                 `json:"title,omitempty"`
	Description string                 `json:"description,omitempty"`
	Timestamp   string                 `json:"timestamp,omitempty"`
	Color       int                    `json:"color,omitempty"`
	Footer      *MessageEmbedFooter    `json:"footer,omitempty"`
	Image       *MessageEmbedImage     `json:"image,omitempty"`
	Thumbnail   *MessageEmbedThumbnail `json:"thumbnail,omitempty"`
	Video       *MessageEmbedVideo     `json:"video,omitempty"`
	Provider    *MessageEmbedProvider  `json:"provider,omitempty"`
	Author      *MessageEmbedAuthor    `json:"author,omitempty"`
	Fields      []*MessageEmbedField   `json:"fields,omitempty"`
}

func (w *Webhook) PostMessage() (err error) {
	log.Println("Sending webhook to discord...")

	var msg []byte
	if len(w.Message.Images) > 0 {
		msg, err = w.executeMultipart(false)
	} else {
		msg, err = w.executeJSON(false)
	}
	if err != nil {
		return
	}

	response, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		return
	}

	log.Printf("Webhook response: %s\n", string(response))

	return
}

func (w *Webhook) URL() string {
	return webhookBaseURL + w.ID + "/" + w.Token
}

func (w *Webhook) executeJSON(wait bool) (response []byte, err error) {
	url := w.URL()
	if wait {
		url += "?wait=true"
	}

	log.Println("Issuing webhook to", url)

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, url, body)
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	log.Println("discord response:", resp.Status)

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return respBody, nil
}

func (w *Webhook) executeMultipart(wait bool) (response []byte, err error) {
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)

	// add image form fields
	err = addImages(mw, w.Params.Images)
	if err != nil {
		return
	}

	// add other content in form field "payload_json"
	if w.Params.Content != "" {
		err = addPayloadJSON(mw, w.Params.Content)
		if err != nil {
			return
		}
	}

	mw.Close()

	url := w.URL()
	if wait {
		url += "?wait=true"
	}

	log.Println("Issuing webhook to", url)

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, url, body)
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", mw.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	log.Println("discord response:", resp.Status)

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return respBody, nil
}

func addPayloadJSON(w *multipart.Writer, content string) error {
	jsonPayload, err := json.Marshal(map[string]string{
		"content": content,
	})
	if err != nil {
		return fmt.Errorf("error marshalling content: %w", err)
	}

	fw, err := w.CreateFormField("payload_json")
	if err != nil {
		return fmt.Errorf("unable to create payload_json field: %w", err)
	}

	_, err = fw.Write(jsonPayload)
	if err != nil {
		return fmt.Errorf("error writing json payload: %w", err)
	}

	return nil
}

func addImages(w *multipart.Writer, images []*plexwebhooks.Thumbnail) error {
	for _, img := range images {
		fw, err := w.CreateFormFile("image", img.Filename)
		if err != nil {
			return fmt.Errorf("error creating image part: %w", err)
		}

		buf := bytes.NewBuffer(img.Data)

		thumb, err := jpeg.Decode(buf)
		if err != nil {
			return fmt.Errorf("error decoding thumbnail: %w", err)
		}

		err = jpeg.Encode(fw, thumb, nil)
		if err != nil {
			return fmt.Errorf("error encoding thumbnail jpeg: %w", err)
		}
	}

	return nil
}
