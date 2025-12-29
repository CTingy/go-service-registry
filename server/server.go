package main

import (
	"fmt"
	"log"
	"net"

	"go-service-registry/pkg/storage"
	pb "go-service-registry/proto"

	"google.golang.org/grpc"
)

func main() {
	const (
		port       = ":50055"
		shardCount = 32
		defaultTTL = 15
	)

	store := storage.NewShardedMap(shardCount)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	registryServer := NewRegistryServer(store, int64(defaultTTL))
	// naming: Register + {service name: Registry} + Server
	pb.RegisterRegistryServer(grpcServer, registryServer)

	fmt.Printf("Registry Server is running on port %s\n", port)
	fmt.Printf("Settings: TTL=%ds, Shards=%d\n", defaultTTL, shardCount)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
