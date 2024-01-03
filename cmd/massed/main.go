// Package main provides the massed buildkit gateway.
package main

import (
	"github.com/moby/buildkit/frontend/gateway/grpcclient"
	"github.com/moby/buildkit/util/appcontext"
	"github.com/moby/buildkit/util/bklog"

	"gitlab.wikimedia.org/dduvall/masse/gateway"
)

func main() {
	err := grpcclient.RunFromEnvironment(appcontext.Context(), gateway.Run)

	if err != nil {
		bklog.L.Errorf("fatal error: %+v", err)
		panic(err)
	}
}
