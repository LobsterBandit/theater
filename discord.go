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
	ID     string
	Token  string
	Params *WebhookParams
}

type WebhookParams struct {
	Content string
	Images  []*plexwebhooks.Thumbnail
	Embeds  []*MessageEmbed
}

type MessageEmbed struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Timestamp   string `json:"timestamp,omitempty"`
	Color       int    `json:"color,omitempty"`
}

func (w *Webhook) PostMessage() (err error) {
	log.Println("Sending webhook to discord...")

	msg, err := w.executeMultipart(false)
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
