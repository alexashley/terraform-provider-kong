package provider

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Kong struct {
	adminApiUrl string
	client      *http.Client
}

func NewKongClient(adminApiUrl string) (*Kong, error) {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	return &Kong{adminApiUrl: adminApiUrl, client: httpClient}, nil
}

func (kong *Kong) createService(service KongService) (string, error) {
	servicesEndpoint := kong.adminApiUrl + "/services"
	response, err := kong.client.PostForm(servicesEndpoint, url.Values{
		"name":     {service.name},
		"protocol": {service.protocol},
		"host":     {service.host},
		"port":     {strconv.Itoa(service.port)},
		"path":     {service.path},
	})

	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	var otherService KongService
	json.Unmarshal([]byte(body), &otherService)

	return otherService.id, nil
}
