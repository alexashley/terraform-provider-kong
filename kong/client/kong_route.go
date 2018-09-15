package client

const routePath = "/routes"

func (kongClient *KongClient) CreateRoute(routeToCreate KongRoute) (*KongRoute, error) {
	var newRoute KongRoute
	err := kongClient.post(routePath, routeToCreate, &newRoute)

	if err != nil {
		return nil, err
	}

	return &newRoute, nil
}

func (kongClient *KongClient) GetRoute(routeId string) (*KongRoute, error) {
	var route KongRoute
	err := kongClient.get(routePath+"/"+routeId, &route)

	if err != nil {
		return nil, err
	}

	return &route, nil
}

func (kongClient *KongClient) DeleteRoute(routeId string) error {
	return kongClient.delete(routePath + "/" + routeId)
}

func (kongClient *KongClient) UpdateRoute(routeToUpdate KongRoute) error {
	return kongClient.put(routePath+"/"+routeToUpdate.Id, routeToUpdate)
}
