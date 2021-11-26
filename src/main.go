package main

import (
	"log"
	"main/handlers"
	"net/http"
	"os"
)

type ServerConfig struct {
	BindAddr string
}

// server config, values here will get overriden by env
var config = ServerConfig{
	BindAddr: "0.0.0.0:2000",
}

func makeMux() *http.ServeMux {
	mux := http.NewServeMux()

	handler := handlers.NewBaseHandler()

	mux.HandleFunc("/v1/lessons", handler.GetLessons)

	return mux
}

func configure(c *ServerConfig) {
	bindAddr, exists := os.LookupEnv("BIND_ADDR")
	if exists {
		if bindAddr == "" {
			log.Fatalln("(cfg) empty bind address supplied, cannot bind")
		}
		c.BindAddr = bindAddr
	} else {
		log.Printf("(cfg) no bind address supplied, defaulting to '%s'", c.BindAddr)
	}
}

func main() {
	log.Printf("Initializing config from environment variables...")

	configure(&config)

	log.Printf("Starting server binded to %s...", config.BindAddr)

	mux := makeMux()
	http.ListenAndServe(config.BindAddr, handlers.RequestLogger(mux))

	log.Println("Server exited. Cleaning up...")
}
