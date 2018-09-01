package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	rbacHeader = "Kong-Admin-Token"
)

type KongConfig struct {
	AdminApiUrl string
	RbacToken   string
}

type KongClient struct {
	Config KongConfig

	client *http.Client
}

func NewKongClient(config KongConfig) (*KongClient, error) {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	return &KongClient{Config: config, client: httpClient}, nil
}

func (kongClient *KongClient) postJson(path string, payload interface{}, responseResource interface{}) error {
	endpoint := kongClient.Config.AdminApiUrl + path

	serializedPayload, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", endpoint, bytes.NewReader(serializedPayload))

	if err != nil {
		return err
	}

	kongClient.addDefaultHeaders(request)
	request.Header.Set("Content-Type", "application/json")

	response, err := kongClient.client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return err
	}

	return json.Unmarshal(body, responseResource)
}

func (kongClient *KongClient) get(path string, responseResource interface{}) error {
	endpoint := kongClient.Config.AdminApiUrl + path
	request, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		return err
	}

	kongClient.addDefaultHeaders(request)

	response, err := kongClient.client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return err
	}

	return json.Unmarshal(body, responseResource)
}

func (kongClient *KongClient) delete(path string) error {
	endpoint := kongClient.Config.AdminApiUrl + path

	request, err := http.NewRequest("DELETE", endpoint, nil)

	if err != nil {
		return err
	}

	_, err = kongClient.client.Do(request)

	if err != nil {
		return err
	}

	return nil
}

func (kongClient *KongClient) addDefaultHeaders(request *http.Request) {
	if kongClient.Config.RbacToken != "" {
		request.Header.Set(rbacHeader, kongClient.Config.RbacToken)
	}

	request.Header.Set("Accept", "application/json")
}
