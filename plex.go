package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"

	"github.com/hekmon/plexwebhooks"
)

var (
	ErrNilMultipartReader    = errors.New("multipart reader cannot be nil")
	ErrMultiplePayloadPart   = errors.New("payload part present more than once")
	ErrMultipleThumbnailPart = errors.New("thumbnail part present more than once")
	ErrPayloadNotFound       = errors.New("payload not found in request")
	ErrUnexpectedFormPart    = errors.New("unexpected form part encountered")
)

func UnexpectedFormPartError(msg string) error {
	return fmt.Errorf("%w: %s", ErrUnexpectedFormPart, msg)
}

type WebhookResult struct {
	RawPayload []byte
	Payload    *plexwebhooks.Payload
	Thumbnail  *plexwebhooks.Thumbnail
}

// ParsePlexWebhook extracts the payload and the thumbnail (if present) from a multipart reader.
func ParsePlexWebhook(mpr *multipart.Reader) (webhook *WebhookResult, err error) { //nolint:cyclop
	webhook = &WebhookResult{}

	if mpr == nil {
		err = ErrNilMultipartReader

		return
	}

	var formPart *multipart.Part
	for formPart, err = mpr.NextPart(); err == nil; formPart, err = mpr.NextPart() {
		switch formPart.FormName() {
		case "payload":
			if webhook.Payload != nil {
				err = ErrMultiplePayloadPart

				return
			}

			buf := new(bytes.Buffer)
			if _, err = buf.ReadFrom(formPart); err != nil {
				err = fmt.Errorf("payload form part read failed: %w", err)

				return
			}

			webhook.RawPayload = buf.Bytes()

			webhook.Payload = new(plexwebhooks.Payload)
			if err = json.Unmarshal(webhook.RawPayload, webhook.Payload); err != nil {
				err = fmt.Errorf("payload JSON decode failed: %w", err)

				return
			}
		case "thumb":
			if webhook.Thumbnail != nil {
				err = ErrMultipleThumbnailPart

				return
			}

			webhook.Thumbnail = &plexwebhooks.Thumbnail{
				Filename: formPart.FileName(),
			}

			if webhook.Thumbnail.Data, err = ioutil.ReadAll(formPart); err != nil {
				err = fmt.Errorf("error while reading thumb form part data: %w", err)

				return
			}
		default:
			err = UnexpectedFormPartError(formPart.FormName())

			return
		}
	}

	if errors.Is(err, io.EOF) {
		err = nil
	}

	if err == nil && webhook.Payload == nil {
		err = ErrPayloadNotFound
	}

	return webhook, err
}
