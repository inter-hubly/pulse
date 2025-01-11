package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/inter-hubly/pilot/hlog"
	"github.com/inter-hubly/pilot/server"
	"github.com/inter-hubly/pulse/internal/infraestructure/express"
)

func main() {
	ctx := context.Background()

	server.FillConfigEnvironment(ctx)

	express.Start(ctx)

	hlog.Info(ctx, "main", fmt.Sprintf("Server start in port %s", server.GetEnvironment().Port))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", server.GetEnvironment().Port), nil); err != nil {
		log.Fatal(err)
	}
}
