package main

import (
	"context"
	"flag"
	"log"

	"github.com/denouche/terraform-provider-awx/awx"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "run the provider with support for debuggers like delve")
	flag.Parse()

	ctx := context.Background()

	// The SDKv2 provider serves the existing resources and data sources; the
	// plugin-framework provider serves what only the framework supports, such
	// as actions. Both share one provider schema and configuration.
	muxServer, err := tf5muxserver.NewMuxServer(ctx,
		awx.Provider().GRPCProvider,
		providerserver.NewProtocol5(awx.NewFrameworkProvider()),
	)
	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf5server.ServeOpt
	if debug {
		serveOpts = append(serveOpts, tf5server.WithManagedDebug())
	}

	err = tf5server.Serve(
		"registry.terraform.io/denouche/awx",
		func() tfprotov5.ProviderServer { return muxServer.ProviderServer() },
		serveOpts...,
	)
	if err != nil {
		log.Fatal(err)
	}
}
