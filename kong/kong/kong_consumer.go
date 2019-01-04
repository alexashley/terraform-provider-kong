package kong

import "strings"

const consumerPath = "/consumers"

func singleConsumer(consumerId string) string {
	return consumerPath + "/" + consumerId
}

func (kongClient *KongClient) CreateConsumer(consumerToCreate *KongConsumer) (*KongConsumer, error) {
	var newConsumer KongConsumer
	err := kongClient.post(consumerPath, consumerToCreate, &newConsumer)

	if err != nil {
		return nil, err
	}

	return &newConsumer, nil
}

func (kongClient *KongClient) GetConsumer(consumerId string) (*KongConsumer, error) {
	var consumer KongConsumer

	err := kongClient.get(singleConsumer(consumerId), &consumer)

	if err != nil {
		return nil, err
	}

	return &consumer, nil
}

func (kongClient *KongClient) DeleteConsumer(consumerId string) error {
	return kongClient.delete(singleConsumer(consumerId))
}

func (kongClient *KongClient) UpdateConsumer(consumerToUpdate *KongConsumer) error {
	return kongClient.put(singleConsumer(consumerToUpdate.Id), consumerToUpdate)
}

func (kongClient *KongClient) GetConsumers() ([]KongConsumer, error) {
	var page KongConsumerPage
	var plugins []KongConsumer

	for {
		url := consumerPath

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
