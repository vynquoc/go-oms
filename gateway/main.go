package main

import (
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	common "github.com/vynquoc/go-oms-common"
)

var (
	addr = common.EnvString("HTTP_ADDR", ":8080")
)

func main() {
	mux := http.NewServeMux()
	handler := NewHandler()
	handler.registerRoutes(mux)

	log.Printf("Starting server at %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal("Failed to create server")
	}
}
