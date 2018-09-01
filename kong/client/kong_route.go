package client

const routePath = "/routes"

func (kongClient *KongClient) CreateRoute(routeToCreate KongRoute) (string, error) {
	var newRoute KongRoute
	err := kongClient.postJson(routePath, routeToCreate, &newRoute)

	if err != nil {
		return "", err
	}

	return newRoute.Id, nil
}

func (kongClient *KongClient) DeleteRoute(routeId string) error {
	return kongClient.delete(routePath + "/" + routeId)
}
