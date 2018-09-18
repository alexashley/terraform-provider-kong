package kong

import (
	"fmt"
	"testing"
)

func TestGetStatus(t *testing.T) {
	status, err := kongClient.GetStatus()

	if err != nil {
		t.Errorf("Failed to get Kong status %s", err.Error())
		return
	}

	fmt.Printf("Database reachable? %t", status.Database.Reachable)
}

func TestGetStatusFailure(t *testing.T) {
	_, err := NewKongClient(KongConfig{
		AdminApiUrl: "http://foo.bar",
	})

	if err == nil {
		t.Errorf("expected an error when initializing Kong with an invalid url")
	}
}
