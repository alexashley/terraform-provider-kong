package kong

import "strings"

const pluginPath = "/plugins"

func (kongClient *KongClient) CreatePlugin(pluginToCreate *KongPlugin) (*KongPlugin, error) {
	var newPlugin KongPlugin

	err := kongClient.post(pluginPath, pluginToCreate, &newPlugin)

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

func (kongClient *KongClient) UpdatePlugin(pluginToUpdate *KongPlugin) error {
	return kongClient.patch(pluginPath+"/"+pluginToUpdate.Id, pluginToUpdate)
}

func (kongClient *KongClient) GetPlugins() ([]KongRoute, error) {
	var page KongPluginsPage
	var plugins []KongRoute

	for {
		url := pluginPath

		if page.Next != "" {
			url = strings.Replace(page.Next, kongClient.Config.AdminApiUrl, "", 1)
		}

		if err := kongClient.get(url, &page); err != nil {
			return nil, err
		}

		plugins = append(plugins, page.Data...)

		if page.Next == "" {
			break
		}
	}

	return plugins, nil
}
