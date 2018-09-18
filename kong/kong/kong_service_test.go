package kong

import (
	"fmt"
	"testing"
)

var kongClient *KongClient

func init() {
	kongClient, _ = NewKongClient(KongConfig{
		AdminApiUrl: "http://localhost:8001",
	})
}

func TestGetAllServices(t *testing.T) {
	services, err := kongClient.GetServices()

	if err != nil {
		t.Errorf("Failed to get all Kong services %s", err.Error())
		return
	}

	fmt.Printf("Got %d services\n", len(services))
}
