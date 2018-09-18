package kong

func (kongClient *KongClient) GetStatus() (*KongStatus, error) {
	var kongStatus KongStatus

	err := kongClient.get("/status", &kongStatus)

	return &kongStatus, err
}
