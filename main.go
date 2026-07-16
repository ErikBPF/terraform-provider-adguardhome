package main

import (
	"context"

	"github.com/ErikBPF/terraform-provider-adguardhome/adguard"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// provider documentation generation
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name adguardhome

func main() {
	providerserver.Serve(context.Background(), adguard.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/ErikBPF/adguardhome",
	})
}
