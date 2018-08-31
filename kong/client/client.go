package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type KongClient struct {
	AdminApiUrl string

	client *http.Client
}

func NewKongClient(adminApiUrl string) (*KongClient, error) {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	return &KongClient{AdminApiUrl: adminApiUrl, client: httpClient}, nil
}

func (kongClient *KongClient) CreateService(serviceToCreate KongService) (string, error) {
	servicePayload := url.Values{
		"name":     {serviceToCreate.Name},
		"protocol": {serviceToCreate.Protocol},
		"host":     {serviceToCreate.Host},
		"port":     {strconv.Itoa(serviceToCreate.Port)},
		"path":     {serviceToCreate.Path},
	}

	var newService KongService
	err := kongClient.postForm("/services", servicePayload, &newService)

	if err != nil {
		return "", err
	}

	return newService.Id, nil
}

func (kongClient *KongClient) postForm(path string, form url.Values, responseResource interface{}) error {
	endpoint := kongClient.AdminApiUrl + path

	request, err := http.NewRequest("POST", endpoint, strings.NewReader(form.Encode()))

	if err != nil {
		return err
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
