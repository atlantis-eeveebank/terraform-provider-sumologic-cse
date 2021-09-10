package main

import (
	sumologiccse "github.com/eeveebank/terraform-provider-sumologic-cse/sumologic-cse"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return sumologiccse.Provider()
		},
	})
}
