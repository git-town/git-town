package bitbucket

import (
	"encoding/json"
	"errors"

	"github.com/mitchellh/mapstructure"
)

type SSHKeys struct {
	c *Client
}

type SSHKey struct {
	Uuid      string `json:"uuid"`
	Label     string `json:"label"`
	Key       string `json:"key"`
	Comment   string `json:"comment"`
	CreatedOm string `json:"created_on"`
}

type SSHKeyRes struct {
	Page     int32
	Pagelen  int32
	MaxDepth int32
	Size     int32
	Items    []SSHKey
}

func decodeSSHKey(response interface{}) (*SSHKey, error) {
	respMap := response.(map[string]interface{})

	if respMap["type"] == "error" {
		return nil, DecodeError(respMap)
	}

	var sshKey = new(SSHKey)
	err := mapstructure.Decode(respMap, sshKey)
	if err != nil {
		return nil, err
	}

	return sshKey, nil
}

func buildSSHKeysBody(opt *SSHKeyOptions) (string, error) {
	body := map[string]interface{}{}
	body["label"] = opt.Label
	body["key"] = opt.Key

	data, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func decodeSSHKeys(keysResponse interface{}) (*SSHKeyRes, error) {
	keysResponseMap, ok := keysResponse.(map[string]interface{})
	if !ok {
		return nil, errors.New("Not a valid format")
	}

	keyArray := keysResponseMap["values"].([]interface{})
	var keys []SSHKey
	for _, keyEntry := range keyArray {
		var key SSHKey
		err := mapstructure.Decode(keyEntry, &key)
		if err == nil {
			keys = append(keys, key)
		}
	}

	page, ok := keysResponseMap["page"].(float64)
	if !ok {
		page = 0
	}

	pagelen, ok := keysResponseMap["pagelen"].(float64)
	if !ok {
		pagelen = 0
	}
	max_depth, ok := keysResponseMap["max_width"].(float64)
	if !ok {
		max_depth = 0
	}
	size, ok := keysResponseMap["size"].(float64)
	if !ok {
		size = 0
	}

	keysResp := &SSHKeyRes{
		Page:     int32(page),
		Pagelen:  int32(pagelen),
		MaxDepth: int32(max_depth),
		Size:     int32(size),
		Items:    keys,
	}
	return keysResp, nil
}

func (sk *SSHKeys) Create(ro *SSHKeyOptions) (*SSHKey, error) {
	data, err := buildSSHKeysBody(ro)
	if err != nil {
		return nil, err
	}
	urlStr := sk.c.requestUrl("/users/%s/ssh-keys", ro.Owner)
	response, err := sk.c.execute("POST", urlStr, data)
	if err != nil {
		return nil, err
	}

	return decodeSSHKey(response)
}

func (sk *SSHKeys) Get(ro *SSHKeyOptions) (*SSHKey, error) {
	urlStr := sk.c.requestUrl("/users/%s/ssh-keys/%s", ro.Owner, ro.Uuid)
	response, err := sk.c.execute("GET", urlStr, "")
	if err != nil {
		return nil, err
	}

	return decodeSSHKey(response)
}

func (sk *SSHKeys) Delete(ro *SSHKeyOptions) (interface{}, error) {
	urlStr := sk.c.requestUrl("/users/%s/ssh-keys/%s", ro.Owner, ro.Uuid)
	return sk.c.execute("DELETE", urlStr, "")
}
