package main

import (
	"context"
	"log"

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

/*
type RegistryServer interface {
	Register(context.Context, *RegisterReq) (*RegisterResp, error)
	Heartbeat(context.Context, *HeartbeatReq) (*Empty, error)
	Discover(context.Context, *DiscoverReq) (*DiscoverResp, error)
	mustEmbedUnimplementedRegistryServer()
}
*/

func (s *RegistryServer) Register(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterResp, error) {
	if req.Endpoint == "" || req.ServiceName == "" {
		log.Println("Register rejected: Missing service name or endpoint.")
		return &pb.RegisterResp{Ttl: s.defaultTTL}, nil
	}

	// write data
	s.store.Register(req.ServiceName, req.Endpoint, s.defaultTTL)
	log.Printf("[Register] Service: %s, Endpoint: %s", req.ServiceName, req.Endpoint)

	// return TTL to client
	return &pb.RegisterResp{Ttl: s.defaultTTL}, nil
}

func (s *RegistryServer) Heartbeat(ctx context.Context, req *pb.HeartbeatReq) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (s *RegistryServer) Discover(ctx context.Context, req *pb.DiscoverReq) (*pb.DiscoverResp, error) {
	return &pb.DiscoverResp{}, nil
}
