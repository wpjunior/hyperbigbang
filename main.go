package main

import (
	"log"
	"net/http"
	"os"

	"github.com/wpjunior/hyperbigbang/handler"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	log.Printf("HyperBigBang is running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, &handler.Handler{}))
}
