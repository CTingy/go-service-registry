package main

import (
	"context"
	"log"

	// ⚠️ 注意：如果你的 go.mod module 名稱不是 go-service-registry，請修改這裡
	"go-service-registry/pkg/storage"
	pb "go-service-registry/proto"
)

// RegistryServer 是我們 gRPC 服務的主體
// 它必須實作 pb.RegistryServer 介面
type RegistryServer struct {
	pb.UnimplementedRegistryServer //這是 gRPC 的規範，必須嵌入這個 struct

	storage    *storage.ShardedMap // 我們的核心引擎
	defaultTTL int64               // 預設的過期時間 (秒)
}

// NewRegistryServer 是一個建構子，用來注入依賴 (Dependency Injection)
func NewRegistryServer(store *storage.ShardedMap, ttl int64) *RegistryServer {
	return &RegistryServer{
		storage:    store,
		defaultTTL: ttl,
	}
}

// Register 實作註冊邏輯
// 當 Client 呼叫 Register 時，這個 function 會被執行
func (s *RegistryServer) Register(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterResp, error) {
	// 簡單的參數檢查
	if req.ServiceName == "" || req.Endpoint == "" {
		log.Println("Register rejected: missing service name or endpoint")
		// 雖然參數錯了，但為了簡化，我們先回傳一個預設 TTL，不報錯
		return &pb.RegisterResp{Ttl: int32(s.defaultTTL)}, nil
	}

	// 呼叫我們的核心引擎寫入資料
	s.storage.Register(req.ServiceName, req.Endpoint, s.defaultTTL)

	log.Printf("[Register] Service: %s, Endpoint: %s", req.ServiceName, req.Endpoint)

	// 回傳 TTL 給 Client
	return &pb.RegisterResp{
		Ttl: int32(s.defaultTTL),
	}, nil
}

// Heartbeat 實作心跳邏輯
func (s *RegistryServer) Heartbeat(ctx context.Context, req *pb.HeartbeatReq) (*pb.Empty, error) {
	// 心跳其實就是「再註冊一次」，更新 LastHeartbeat 時間
	// 這裡直接複用 Register 的底層邏輯
	s.storage.Register(req.ServiceName, req.Endpoint, s.defaultTTL)

	// Log 太吵的話可以註解掉這行
	// log.Printf("[Heartbeat] Service: %s", req.ServiceName)

	return &pb.Empty{}, nil
}

// Discover 實作服務發現邏輯
func (s *RegistryServer) Discover(ctx context.Context, req *pb.DiscoverReq) (*pb.DiscoverResp, error) {
	// 從核心引擎讀取 IP 列表
	endpoints := s.storage.GetEndpoints(req.ServiceName)

	log.Printf("[Discover] Query: %s, Found: %d endpoints", req.ServiceName, len(endpoints))

	return &pb.DiscoverResp{
		Endpoints: endpoints,
	}, nil
}
