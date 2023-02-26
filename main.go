package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/winebarrel/terraform-provider-redash/redash"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name redash

func main() {
	debug := flag.Bool("debug", false, "debug mode")
	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: redash.Provider,
		ProviderAddr: "registry.terraform.io/winebarrel/redash",
		Debug:        *debug,
	}

	plugin.Serve(opts)
}
