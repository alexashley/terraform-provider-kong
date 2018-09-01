package client

import (
	"net/url"
	"strconv"
)

const servicesPath = "/services"

func (kongClient *KongClient) CreateService(serviceToCreate KongService) (string, error) {
	servicePayload := createNewServiceFormBody(serviceToCreate)

	var newService KongService
	err := kongClient.postForm(servicesPath, servicePayload, &newService)

	if err != nil {
		return "", err
	}

	return newService.Id, nil
}

func (kongClient *KongClient) DeleteService(serviceId string) error {
	return kongClient.delete(servicePath(serviceId))
}

func (kongClient *KongClient) GetService(serviceId string) (*KongService, error) {
	var service KongService
	err := kongClient.get(servicePath(serviceId), &service)

	if err != nil {
		return nil, err
	}

	service.Url = service.Protocol + service.Host + strconv.Itoa(service.Port) + service.Path

	return &service, nil
}

func createNewServiceFormBody(serviceToCreate KongService) url.Values {
	if serviceToCreate.Url != "" {
		return url.Values{
			"name": {serviceToCreate.Name},
			"url":  {serviceToCreate.Url},
		}
	}

	return url.Values{
		"name":     {serviceToCreate.Name},
		"protocol": {serviceToCreate.Protocol},
		"host":     {serviceToCreate.Host},
		"port":     {strconv.Itoa(serviceToCreate.Port)},
		"path":     {serviceToCreate.Path},
	}
}

func servicePath(serviceId string) string {
	return servicesPath + "/" + serviceId
}
