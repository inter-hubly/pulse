package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/inter-hubly/pilot/server"
	"github.com/inter-hubly/pulse/internal/infraestructure/express"
)

func main() {
	server.FillConfigEnvironment()
	express.Start()

	if err := http.ListenAndServe(fmt.Sprintf(":%s", server.GetEnvironment().Port), nil); err != nil {
		log.Fatal(err)
	}
}
