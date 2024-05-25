package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/stripe/stripe-go/v78"
	common "github.com/vynquoc/go-oms-common"
	"github.com/vynquoc/go-oms-common/broker"
	"github.com/vynquoc/go-oms-common/discovery"
	"github.com/vynquoc/go-oms-common/discovery/consul"
	stripeProcessor "github.com/vynquoc/go-oms-payments/processor/stripe"
	"google.golang.org/grpc"
)

var (
	serviceName  = "payments"
	grpcAddr     = common.EnvString("GRPC_ADDR", "localhost:3003")
	consulAddr   = common.EnvString("CONSUL_ADDR", "localhost:8500")
	amqpUser     = common.EnvString("RABBITMQ_USER", "guest")
	amqpPassword = common.EnvString("RABBITMQ_PASSWORD", "guest")
	amqpHost     = common.EnvString("RABBITMQ_HOST", "localhost")
	amqpPort     = common.EnvString("RABBITMQ_PORT", "5672")
	stripeKey    = common.EnvString("STRIPE_KEY", "sk_test_51PFVb72KofBMjvPrdSHUUhmLzskCch287HbIfFs5e9bKDwN6Xn7Y8JQoFExsgSr3n523CsJ9H9XmJJxIsieDm3Nc003XvMAYOS")
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

	//stripe setup
	stripe.Key = stripeKey

	// broker connection
	ch, close := broker.Connect(amqpUser, amqpPassword, amqpHost, amqpPort)

	defer func() {
		close()
		ch.Close()
	}()
	stripeProcessor := stripeProcessor.NewProcessor()
	svc := NewService(stripeProcessor)

	amqpConsumer := NewConsumer(svc)

	go amqpConsumer.Listen(ch)

	grpcServer := grpc.NewServer()
	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen %v", grpcAddr)
	}
	defer l.Close()

	log.Println("GRPC Order service started at:", grpcAddr)

	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
