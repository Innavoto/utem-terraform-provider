package main

import (
	"context"
	"log"

	"github.com/Innavoto/terraform-provider-utem/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var version = "0.1.0"

func main() {
	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/Innavoto/utem",
	}
	err := providerserver.Serve(context.Background(), provider.New(version), opts)
	if err != nil {
		log.Fatal(err)
	}
}
