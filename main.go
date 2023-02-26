package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/winebarrel/terraform-provider-redash/redash"
)

func main() {
	debug := flag.Bool("debug", false, "debug mode")
	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: redash.Provider,
	}

	if *debug {
		err := plugin.Debug(context.Background(), "registry.terraform.io/winebarrel/redash", opts)

		if err != nil {
			log.Fatal(err)
		}
	} else {
		plugin.Serve(opts)
	}
}
