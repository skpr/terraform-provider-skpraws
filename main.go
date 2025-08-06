package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/skpr/terraform-provider-skpraws/internal/provider"
)

var (
	// set by goreleaser.
	version string = "dev"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		// TODO: Update this string with the published name of your provider.
		// Also update the tfplugindocs generate command to either remove the
		// -provider-name flag or set its value to the updated provider name.
		Address: "registry.terraform.io/skpr/skpraws",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.NewSkprAwsProvider(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
