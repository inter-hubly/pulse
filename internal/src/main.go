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

	log.Println("HTTP server started on :8082")
	if err := http.ListenAndServe(fmt.Sprintf(":%d", server.GetEnvironment().Port), nil); err != nil {
		log.Fatal(err)
	}
}
