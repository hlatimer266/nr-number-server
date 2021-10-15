package main

import (
	"context"

	"github.com/hlatimer266/nr-number-server/internal/server"
)

const port = ":4000"

func main() {
	// start up application listening on port 4000
	ctx := context.Background()
	server.Run(port, ctx)
}
