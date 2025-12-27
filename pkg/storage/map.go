package storage

import (
	"hash/fnv"
	"sync"
	"time"
)

type InstanceInfo struct {
	ServiceName   string
	Endpoint      string // e.g 127.0.0.8.1:8000
	LastHeartbeat int64  // unix timestamp
}

type shard struct {
	// lock embedding
	sync.RWMutex

	// {ServiceName: {Endpoint: InstanceInfo}}
	m map[string]map[string]*InstanceInfo
}

// shard array
type ShardedMap struct {
	shards []*shard
	count  int
}

func NewShardedMap(shardCount int) *ShardedMap {
	s := &ShardedMap{
		count:  shardCount,
		shards: make([]*shard, shardCount),
	}

	// init shards
	for i := 0; i < shardCount; i++ {
		s.shards[i] = &shard{
			m: make(map[string]map[string]*InstanceInfo),
		}
	}
	return s
}

func (s *ShardedMap) getShardIndex(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	hashValue := h.Sum32()
	return int(hashValue) % s.count
}

func (s *ShardedMap) getShard(key string) *shard {
	return s.shards[s.getShardIndex(key)]
}

// Register a new service (write)
func (s *ShardedMap) Register(serviceName, endpoint string, ttl int64) {

	sh := s.getShard(serviceName)

	sh.Lock()
	defer sh.Unlock()

	if _, ok := sh.m[serviceName]; !ok {
		sh.m[serviceName] = make(map[string]*InstanceInfo)
	}

	sh.m[serviceName][endpoint] = &InstanceInfo{
		ServiceName: serviceName,
		Endpoint: endpoint,
		LastHeartbeat: time.Now().Unix(),
	}
}

// GetEndpoints (read)
func (s *ShardedMap) GetEndpoints(serviceName string) []string {
	sh := s.getShard(serviceName)
	sh.RLock()
	defer sh.RUnlock()

	// init a String Slice (empty slice. cap = 0)
	endpoints := make([]string, 0)

	if instances, ok := sh.m[serviceName]; ok {
		for ep := range instances {
			endpoints = append(endpoints, ep)
		}
	}
	return endpoints
}

// ReapLoop goroutine: clear timeout service instance
func (s *ShardedMap) ReapLoop(ttlSeconds int64) {
	for _, item := range s.shards {
		go func(sh *shard) {
			ticker := time.NewTicker(1 * time.Second)

			for range ticker.C {
				// write lock for a single shard
				sh.Lock()

				now := time.Now().Unix()
				for serviceName, infos := range sh.m {
					for endpoint, instanceInfo := range infos {
						if now-instanceInfo.LastHeartbeat > ttlSeconds {
							delete(infos, endpoint)
						}
					}
					if len(infos) == 0 {
						delete(sh.m, serviceName)
					}
				}
				sh.Unlock()
			}
		}(item)
	}
}
