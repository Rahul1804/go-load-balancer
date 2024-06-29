package main

import (
	"go-load-balancer/internal/config"
	"go-load-balancer/internal/handler"
	"go-load-balancer/internal/server"
	"go-load-balancer/pkg/log"
)

func main() {
	log.Init()

	config, err := config.LoadConfig("config.json")
	if err != nil {
		log.Error.Fatalf("Error loading config: %v", err)
	}

	lb := handler.NewLoadBalancer(config.Servers, 5)
	server.StartServer(lb, "8080")
}
