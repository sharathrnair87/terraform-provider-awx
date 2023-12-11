package main 

import (
  "context" 
  "github.com/hashicorp/terraform-plugin-framework/providerserver"
  "awx/awx"
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
      Debug: false,
    })

  if err != nil {
    log.Fatal(err.Error())
  }
}

