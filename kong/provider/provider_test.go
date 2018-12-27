package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"testing"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = KongProvider()
	testAccProviders = map[string]terraform.ResourceProvider{
		"kong": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := KongProvider().InternalValidate(); err != nil {
		t.Fatalf("err: %s ", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if adminApiUrl := os.Getenv("KONG_ADMIN_API_URL"); adminApiUrl == "" {
		t.Fatal("KONG_ADMIN_API_URL must be set for acceptance tests")
	}

	err := testAccProvider.Configure(terraform.NewResourceConfig(nil))
	if err != nil {
		t.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	exitCode := m.Run()
	if exitCode == 0 && testing.CoverMode() != "" {
		c := testing.Coverage()
		if c < 0.59 {
			fmt.Println("Coverage failure: ", c)
			exitCode = -1
		}
	}
	os.Exit(exitCode)
}
