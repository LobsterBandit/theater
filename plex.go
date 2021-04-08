package main

import (
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
	err        error
}

// Extract extracts the payload and the thumbnail (if present) from a multipart reader.
func ParsePlexWebhook(mpr *multipart.Reader) (webhook *WebhookResult) {
	webhook = &WebhookResult{}

	if mpr == nil {
		webhook.err = ErrNilMultipartReader

		return
	}

	for formPart, err := mpr.NextPart(); err == nil; formPart, err = mpr.NextPart() {
		switch formPart.FormName() {
		case "payload":
			if webhook.Payload != nil {
				webhook.err = ErrMultiplePayloadPart

				return
			}

			var someBytes []byte
			if _, err := formPart.Read(someBytes); err != nil {
				webhook.err = fmt.Errorf("payload form part read failed: %w", err)

				return
			}

			// var decodedPayload interface{}
			// if err := json.NewDecoder(formPart).Decode(&decodedPayload); err != nil {
			// 	webhook.err = fmt.Errorf("payload form part read failed: %w", err)

			// 	return
			// }

			// fmt.Println(decodedPayload)

			// webhook.RawPayload = []byte(decodedPayload)

			// webhook.Payload = new(plexwebhooks.Payload)
			// if err := json.Unmarshal(webhook.RawPayload, webhook.Payload); err != nil {
			// 	webhook.err = fmt.Errorf("payload JSON decode failed: %w", err)

			// 	return
			// }
		case "thumb":
			if webhook.Thumbnail != nil {
				webhook.err = ErrMultipleThumbnailPart

				return
			}

			webhook.Thumbnail = &plexwebhooks.Thumbnail{
				Filename: formPart.FileName(),
			}

			if webhook.Thumbnail.Data, err = ioutil.ReadAll(formPart); err != nil {
				webhook.err = fmt.Errorf("error while reading thumb form part data: %w", err)

				return
			}
		default:
			webhook.err = UnexpectedFormPartError(formPart.FormName())

			return
		}
	}

	if errors.Is(webhook.err, io.EOF) {
		webhook.err = nil
	}

	if webhook.err == nil && webhook.Payload == nil {
		webhook.err = ErrPayloadNotFound
	}

	return
}
