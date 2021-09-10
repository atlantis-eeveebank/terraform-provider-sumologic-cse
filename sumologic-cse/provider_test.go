package sumologic_cse

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
	"testing"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"sumologiccse": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = Provider()
}

func testAccPreCheck(t *testing.T) {
	if err := os.Getenv("SUMOLOGIC_CSE_API_KEY"); err == "" {
		t.Fatal("SUMOLOGIC_CSE_API_KEY must be set for acceptance tests")
	}
	if err := os.Getenv("SUMOLOGIC_CSE_HOST"); err == "" {
		t.Fatal("SUMOLOGIC_CSE_HOST must be set for acceptance tests")
	}
}
