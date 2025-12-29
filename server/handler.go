package main

import (
	"go-service-registry/pkg/storage"
	pb "go-service-registry/proto"
)

type RegistryServer struct {
	pb.UnimplementedRegistryServer
	store      *storage.ShardedMap
	defaultTTL int64
}

func NewRegistryServer(store *storage.ShardedMap, ttl int64) *RegistryServer {
	return &RegistryServer{
		store:      store,
		defaultTTL: ttl,
	}
}
