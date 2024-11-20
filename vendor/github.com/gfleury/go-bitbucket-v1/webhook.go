package bitbucketv1

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// WebHookRepoPush contains payload to use while reading handling webhooks from bitbucket
type WebHookRepoPush struct {
	Actor      Actor      `json:"actor"`
	Repository Repository `json:"repository"`
	Push       struct {
		Changes []Change `json:"changes"`
	} `json:"push"`
}

// Actor contains the actor of reported changes from a webhook
type Actor struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}

// Change contains changes reported by webhooks
type Change struct {
	Created bool        `json:"created"`
	Closed  bool        `json:"closed"`
	Old     interface{} `json:"old"`
	New     struct {
		Type   string `json:"type"`
		Name   string `json:"name"`
		Target struct {
			Type string `json:"type"`
			Hash string `json:"hash"`
		} `json:"target"`
	} `json:"new"`
}

func TriggerRepoPush(webHookURL string, webHook WebHookRepoPush) (*http.Response, error) {
	payLoad, err := json.Marshal(webHook)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", webHookURL, bytes.NewBuffer(payLoad))
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Event-Key", "repo:push")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("User-Agent", "Bitbucket version: 5.12.0, Post webhook plugin version: 1.6.3")
	req.Header.Set("X-Bitbucket-Type", "server")

	client := &http.Client{}
	return client.Do(req)
}
