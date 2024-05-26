package main

import (
	"context"
	"log"
	"net"
	"time"

	_ "github.com/joho/godotenv/autoload"
	common "github.com/vynquoc/go-oms-common"
	"github.com/vynquoc/go-oms-common/broker"
	"github.com/vynquoc/go-oms-common/discovery"
	"github.com/vynquoc/go-oms-common/discovery/consul"
	"google.golang.org/grpc"
)

var (
	serviceName  = "orders"
	grpcAddr     = common.EnvString("GRPC_ADDR", "localhost:3002")
	consulAddr   = common.EnvString("CONSUL_ADDR", "localhost:8500")
	amqpUser     = common.EnvString("RABBITMQ_USER", "guest")
	amqpPassword = common.EnvString("RABBITMQ_PASSWORD", "guest")
	amqpHost     = common.EnvString("RABBITMQ_HOST", "localhost")
	amqpPort     = common.EnvString("RABBITMQ_PORT", "5672")
)

func main() {
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
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

	ch, close := broker.Connect(amqpUser, amqpPassword, amqpHost, amqpPort)

	defer func() {
		close()
		ch.Close()
	}()

	grpcServer := grpc.NewServer()
	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen %v", grpcAddr)
	}
	defer l.Close()

	store := NewStore()
	svc := NewService(store)

	NewGRPCHandler(grpcServer, svc, ch)

	log.Println("GRPC Order service started at:", grpcAddr)

	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
