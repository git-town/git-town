package bitbucket

import (
	"encoding/json"
	"errors"

	"github.com/mitchellh/mapstructure"
)

type DeployKeys struct {
	c *Client
}

type DeployKey struct {
	Id      int    `json:"id"`
	Label   string `json:"label"`
	Key     string `json:"key"`
	Comment string `json:"comment"`
}

type DeployKeysRes struct {
	Page     int32
	Pagelen  int32
	MaxDepth int32
	Size     int32
	Items    []DeployKey
}

func decodeDeployKey(response interface{}) (*DeployKey, error) {
	respMap := response.(map[string]interface{})

	if respMap["type"] == "error" {
		return nil, DecodeError(respMap)
	}

	var deployKey = new(DeployKey)
	err := mapstructure.Decode(respMap, deployKey)
	if err != nil {
		return nil, err
	}

	return deployKey, nil
}

func decodeDeployKeys(deployKeysResponse interface{}) (*DeployKeysRes, error) {
	deployKeysResponseMap, ok := deployKeysResponse.(map[string]interface{})
	if !ok {
		return nil, errors.New("not a valid format")
	}

	repoArray := deployKeysResponseMap["values"].([]interface{})
	var deployKeys []DeployKey
	for _, deployKeyEntry := range repoArray {
		var deployKey DeployKey
		err := mapstructure.Decode(deployKeyEntry, &deployKey)
		if err == nil {
			deployKeys = append(deployKeys, deployKey)
		}
	}

	page, ok := deployKeysResponseMap["page"].(float64)
	if !ok {
		page = 0
	}

	pagelen, ok := deployKeysResponseMap["pagelen"].(float64)
	if !ok {
		pagelen = 0
	}
	maxDepth, ok := deployKeysResponseMap["max_width"].(float64)
	if !ok {
		maxDepth = 0
	}
	size, ok := deployKeysResponseMap["size"].(float64)
	if !ok {
		size = 0
	}

	repositories := DeployKeysRes{
		Page:     int32(page),
		Pagelen:  int32(pagelen),
		MaxDepth: int32(maxDepth),
		Size:     int32(size),
		Items:    deployKeys,
	}
	return &repositories, nil
}

func buildDeployKeysBody(opt *DeployKeyOptions) (string, error) {
	body := map[string]interface{}{}
	body["label"] = opt.Label
	body["key"] = opt.Key

	data, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (dk *DeployKeys) Create(opt *DeployKeyOptions) (*DeployKey, error) {
	data, err := buildDeployKeysBody(opt)
	if err != nil {
		return nil, err
	}
	urlStr := dk.c.requestUrl("/repositories/%s/%s/deploy-keys", opt.Owner, opt.RepoSlug)
	response, err := dk.c.executeWithContext("POST", urlStr, data, opt.ctx)
	if err != nil {
		return nil, err
	}

	return decodeDeployKey(response)
}

func (dk *DeployKeys) Get(opt *DeployKeyOptions) (*DeployKey, error) {
	urlStr := dk.c.requestUrl("/repositories/%s/%s/deploy-keys/%d", opt.Owner, opt.RepoSlug, opt.Id)
	response, err := dk.c.execute("GET", urlStr, "")
	if err != nil {
		return nil, err
	}

	return decodeDeployKey(response)
}

func (dk *DeployKeys) Delete(opt *DeployKeyOptions) (interface{}, error) {
	urlStr := dk.c.requestUrl("/repositories/%s/%s/deploy-keys/%d", opt.Owner, opt.RepoSlug, opt.Id)
	return dk.c.execute("DELETE", urlStr, "")
}

func (dk *DeployKeys) List(opt *DeployKeyOptions) (*DeployKeysRes, error) {
	urlStr := dk.c.requestUrl("/repositories/%s/%s/deploy-keys", opt.Owner, opt.RepoSlug)
	response, err := dk.c.execute("GET", urlStr, "")
	if err != nil {
		return nil, err
	}

	return decodeDeployKeys(response)
}
