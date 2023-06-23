package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/tidal-engineering/terraform-provider-spinnaker/spinnaker"
)

func main() {

	plugin.Serve(
		&plugin.ServeOpts{
			ProviderFunc: spinnaker.New,
			ProviderAddr: "app.terraform.io/SLUS-DCP/spinnaker",
		},
	)
}
