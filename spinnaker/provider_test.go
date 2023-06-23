package spinnaker

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = New()
	testAccProviders = map[string]*schema.Provider{
		"spinnaker": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("GATE_URL") == "" {
		t.Fatal("GATE_URL must be set for acceptance tests")
	}
	err := testAccProvider.Configure(nil, terraform.NewResourceConfigRaw(nil))
	if err != nil {
		t.Fatalf("err: %v", err)
	}
}

func TestProvider(t *testing.T) {
	if err := New().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = New()
}
