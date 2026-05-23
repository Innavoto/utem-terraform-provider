package resources_test

import (
	"github.com/Innavoto/utem-terraform-provider/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories is used by all acceptance tests to configure
// the provider under test. It uses the real provider implementation via the
// Protocol V6 plugin server.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"utem": providerserver.NewProtocol6WithError(provider.New("test")()),
}

// testAccProviderConfig is shared provider configuration for acceptance tests.
// It requires UTEM_API_KEY to be set in the environment for authentication.
const testAccProviderConfig = `
provider "utem" {
  base_url  = "https://utem.innavoto.com"
  tenant_id = "1"
}
`
