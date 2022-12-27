package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type driver struct {
	client *resty.Client
}

type EventsDriver interface {
	MakeEventRequest(context.Context, *Event, string) error
}

func NewEventsHTTPDriver(webhook string, secret string) EventsDriver {
	httpClient := resty.New().
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		EnableTrace().
		SetHeaders(map[string]string{
			"X-Gitlab-Token": secret,
			"Content-Type":   "application/json",
			"Mock":           "true",
		}).
		SetBaseURL(webhook).
		SetRetryCount(3).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(20 * time.Second).
		SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
			return 0, errors.New("quota exceeded")
		})

	return &driver{
		httpClient,
	}
}

func (d *driver) MakeEventRequest(context context.Context, event *Event, eventType string) error {
	// output, _ := json.Marshal(event)
	// log.Info(string(output))
	// return nil
	resp, err := d.client.R().
		SetHeader("X-Gitlab-Event", eventType).
		SetBody(event).
		Post("/")

	if status := resp.StatusCode(); status >= 400 {
		return errors.New(fmt.Sprintf("Error received from server: %s", resp.Status()))
	}

	return err
}
