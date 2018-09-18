package kong

import (
	"strconv"
	"strings"
)

const servicesPath = "/services"

func (kongClient *KongClient) CreateService(serviceToCreate *KongService) (*KongService, error) {
	var newService KongService
	err := kongClient.post(servicesPath, serviceToCreate, &newService)

	if err != nil {
		return nil, err
	}

	return &newService, nil
}

func (kongClient *KongClient) UpdateService(serviceToUpdate *KongService) error {
	return kongClient.put(servicePath(serviceToUpdate.Id), serviceToUpdate)
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

	service.Url = service.getUrlFromParts()

	return &service, nil
}

func (kongClient *KongClient) GetServices() ([]KongService, error) {
	var servicesResult KongServicesPage
	var allServices = make([]KongService, 0)

	for {
		url := servicesPath

		if servicesResult.Next != "" {
			url = strings.Replace(servicesResult.Next, kongClient.Config.AdminApiUrl, "", 1)
		}

		if err := kongClient.get(url, &servicesResult); err != nil {
			return nil, err
		}

		allServices = append(allServices, servicesResult.Data...)

		if servicesResult.Next == "" {
			break
		}
	}

	return allServices, nil
}

func servicePath(serviceId string) string {
	return servicesPath + "/" + serviceId
}

func (kongService *KongService) getUrlFromParts() string {
	port := kongService.Port

	protocol := kongService.Protocol
	url := protocol + "://" + kongService.Host

	if shouldIncludePortInUrl(protocol, port) {
		url += ":" + strconv.Itoa(port)
	}

	return url + kongService.Path
}

func shouldIncludePortInUrl(protocol string, port int) bool {
	return !(protocol == "https" && port == 443 || protocol == "http" && port == 80)
}
