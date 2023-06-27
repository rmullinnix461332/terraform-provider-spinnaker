package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/rmullinnix461332/terraform-provider-spinnaker/spinnaker"
)

func main() {

	plugin.Serve(
		&plugin.ServeOpts{
			ProviderFunc: spinnaker.New,
			ProviderAddr: "app.terraform.io/SLUS-DCP/spinnaker",
		},
	)
}
