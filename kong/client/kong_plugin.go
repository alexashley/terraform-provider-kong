package client

const pluginPath = "/plugins"

func (kongClient *KongClient) CreatePlugin(pluginToCreate KongPlugin) (string, error) {
	var newPlugin KongPlugin

	err := kongClient.postJson(pluginPath, pluginToCreate, &newPlugin)

	if err != nil {
		return "", err
	}

	return newPlugin.Id, nil
}

func (kongClient *KongClient) DeletePlugin(pluginId string) error {
	return kongClient.delete(pluginPath + "/" + pluginId)
}
