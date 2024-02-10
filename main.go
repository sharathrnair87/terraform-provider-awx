package main

import (
	"awx/awx"
	"context"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"log"
)

var (
	version string = "dev"
)

func main() {
	err := providerserver.Serve(
		context.Background(),
		awx.New(version),
		//providerserver.ServeOpts{})
		providerserver.ServeOpts{
			Address: "registry.terraform.io/sharathrnair87/awx",
			Debug:   false,
		})

	if err != nil {
		log.Fatal(err.Error())
	}
}
