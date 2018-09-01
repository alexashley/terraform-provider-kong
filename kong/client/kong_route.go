package client

const routePath = "/routes"

func (kongClient *KongClient) CreateRoute(routeToCreate KongRoute) (*KongRoute, error) {
	var newRoute KongRoute
	err := kongClient.postJson(routePath, routeToCreate, &newRoute)

	if err != nil {
		return nil, err
	}

	return &newRoute, nil
}

func (kongClient *KongClient) DeleteRoute(routeId string) error {
	return kongClient.delete(routePath + "/" + routeId)
}
