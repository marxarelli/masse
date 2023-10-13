// Package main provides the phytond buildkit gateway.
package main

import (
	"log"

	"github.com/moby/buildkit/frontend/gateway/grpcclient"
	"github.com/moby/buildkit/util/appcontext"

	"gitlab.wikimedia.org/dduvall/phyton/gateway"
)

func main() {
	err := grpcclient.RunFromEnvironment(appcontext.Context(), gateway.Run)

	if err != nil {
		log.Panicf("fatal error:\n%v", err)
	}
}
