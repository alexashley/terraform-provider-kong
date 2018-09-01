package client

import (
	"strconv"
)

const servicesPath = "/services"

func (kongClient *KongClient) CreateService(serviceToCreate KongService) (string, error) {
	var newService KongService
	err := kongClient.postJson(servicesPath, serviceToCreate, &newService)

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

	service.Url = service.protocol + service.host + strconv.Itoa(service.port) + service.path

	return &service, nil
}

func servicePath(serviceId string) string {
	return servicesPath + "/" + serviceId
}
