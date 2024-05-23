package main

import (
	"context"
	"log"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
	common "github.com/vynquoc/go-oms-common"
	"github.com/vynquoc/go-oms-common/discovery"
	"github.com/vynquoc/go-oms-common/discovery/consul"
	"github.com/vynquoc/go-oms-gateway/gateway"
)

var (
	serviceName = "gateway"
	addr        = common.EnvString("HTTP_ADDR", ":8080")
	consulAddr  = common.EnvString("CONSUL_ADDR", "localhost:8500")
)

func main() {
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, addr); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Fatal("failed to healthcheck")
			}
			time.Sleep(time.Second * 1)
		}
	}()

	defer registry.DeRegister(ctx, instanceID, serviceName)

	mux := http.NewServeMux()
	ordersGateway := gateway.NewGRPCGateway(registry)
	handler := NewHandler(ordersGateway)
	handler.registerRoutes(mux)

	log.Printf("Starting server at %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal("Failed to create server")
	}
}
