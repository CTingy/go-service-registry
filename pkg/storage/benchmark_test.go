package storage

import (
	"fmt"
	"math/rand"
	"testing"
)

// Test writing (service boot up / heartbeat)
func BenchmarkRegister(b *testing.B) {
	// init ShardMap (32 shards)
	s := NewShardedMap(32)

	b.ResetTimer()
	// default testing goroutine number: GOMAXPROCS (equals to the cpu core number)
	b.RunParallel(func(pb *testing.PB) {
		id := rand.Int()
		counter := 0
		for pb.Next() {
			// give diff service names
			serviceName := fmt.Sprintf("service-%d", id%100)
			endpoint := fmt.Sprintf("10.0.0.%d:%d", id%255, counter)

			s.Register(serviceName, endpoint, 10)
			counter++
		}
	})
}

func BenchmarkDiscover(b *testing.B) {
	s := NewShardedMap(32)

	// Register 100 services. Each service has 10 nodes
	for i := 0; i < 100; i++ {
		serviceName := fmt.Sprintf("service-%d", i)
		for j := 0; j < 10; j++ {
			endpoint := fmt.Sprintf("10.0.0.%d:%d", j, i)
			s.Register(serviceName, endpoint, 10)
		}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		id := rand.Int()
		for pb.Next() {
			serviceName := fmt.Sprintf("service-%d", id%100)
			_ = s.GetEndpoints(serviceName)
		}
	})
}
