package storage

import (
	"hash/fnv"
	"sync"
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
