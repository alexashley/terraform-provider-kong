package client

const pluginPath = "/plugins"

func (kongClient *KongClient) CreatePlugin(pluginToCreate KongPlugin) (*KongPlugin, error) {
	var newPlugin KongPlugin

	err := kongClient.postJson(pluginPath, pluginToCreate, &newPlugin)

	if err != nil {
		return nil, err
	}

	return &newPlugin, nil
}

func (kongClient *KongClient) DeletePlugin(pluginId string) error {
	return kongClient.delete(pluginPath + "/" + pluginId)
}

func (kongClient *KongClient) GetPlugin(pluginId string) (*KongPlugin, error) {
	var kongPlugin KongPlugin
	err := kongClient.get(pluginPath+"/"+pluginId, &kongPlugin)

	if err != nil {
		return nil, err
	}

	return &kongPlugin, nil
}
