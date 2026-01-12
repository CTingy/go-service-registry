package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	// build non-encryption connection
	pb "go-service-registry/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	// 1. set parameters
	// flag.String(param name, default value, help text)
	// command: go run client/main.go -name=payment-service -port=8081
	serviceName := flag.String("name", "auth-service", "Service Name")
	port := flag.String("port", "8080", "Service port")
	flag.Parse()

	// construct client's address
	addr := fmt.Sprintf("127.0.0.1:%s", *port)

	// 2. connect to Registry Server
	target := "127.0.0.1:50055"
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Fail to new client: %v", err)
	}
	defer conn.Close()

	// func NewRegistryClient(cc grpc.ClientConnInterface) RegistryClient
	client := pb.NewRegistryClient(conn)

	// 3. Register
	log.Printf("Registering %s at %s...", *serviceName, addr)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.Register(ctx, &pb.RegisterReq{
		ServiceName: *serviceName,
		Endpoint:    addr,
	})
	if err != nil {
		if status.Code(err) == codes.DeadlineExceeded {
			// grpc status code: 4
			log.Println("Register Deadline Exceeded")
		} else {
			log.Printf("Register failed: %v", err)
		}
		return
	}

	serverTTL := time.Duration(resp.Ttl)
	if serverTTL <= 0 {
		serverTTL = 15
	}
	log.Printf("Register Success! Server TTL: %d seconds", serverTTL)

	// 4. Heartbeat
	go func() {
		// shorter heartbeat interval than ttl. Best practice: ttl/3
		// TODO: reset timer if get a new serverTTL
		ticker := time.NewTicker(serverTTL * time.Second / 3)
		for range ticker.C {
			log.Printf("Sending heartbeat for %s...", *serviceName)

			// TODO: retry if err
			_, err := client.Heartbeat(context.Background(), &pb.HeartbeatReq{
				ServiceName: *serviceName,
				Endpoint:    addr,
			})

			if err != nil {
				log.Printf("Heartbeat failed: %v", err)
			}
		}
	}()

}
