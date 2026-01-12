# Go High-Performance Service Registry

A lightweight, high-concurrency **Service Registry and Discovery** server built with **Go** and **gRPC**.
The system uses a **sharded in-memory architecture** to minimize lock contention and achieve high throughput
under mixed read/write workloads.


## Key Features

- **High-Performance Storage**
  - Implements a **sharded map (32 shards)** to significantly reduce `sync.RWMutex` contention
  - Designed as an alternative to a single global lock or standard Go map

- **gRPC & Protobuf**
  - Efficient binary serialization for low-latency, high-throughput communication

- **Service Governance**
  - Server-side **TTL (Time-To-Live)** enforcement
  - Active **heartbeat monitoring** to automatically evict stale services

- **Thread-Safe by Design**
  - Fully concurrency-safe
  - Optimized for read-heavy and mixed read/write access patterns


## Architecture

The core engine uses a **sharding strategy** (conceptually similar to Java’s `ConcurrentHashMap`)
to partition storage into multiple independent buckets based on the **FNV-1a hash** of the service name.

Each shard owns its own lock, allowing concurrent operations to proceed without contending on a global mutex.

```mermaid
graph TB
  Client -->|Register / Discover| GRPC[gRPC Server]
  GRPC -->|hash service name| Selector[Shard Selector]

  Selector -->|idx = 0| Shard0[Shard 0 - RWMutex]
  Selector -->|idx = 1| Shard1[Shard 1 - RWMutex]
  Selector -->|...| ShardN[Shard N - RWMutex]
  Selector -->|idx = 31| Shard31[Shard 31 - RWMutex]
````


## Benchmark Results


Benchmarks were run on an **Apple M1 (8-core)** machine.
The in-memory storage engine sustains over **4 million operations per second**.

```bash
goos: darwin
goarch: arm64
pkg: go-service-registry/pkg/storage
cpu: Apple M1

BenchmarkRegister-8      5993856    239.6 ns/op   (~4.1M ops/sec)
BenchmarkDiscover-8      4779848    246.2 ns/op   (~4.0M ops/sec)
```


## Quick Start

### Prerequisites

* Go **1.21+**
* `protoc` (Protocol Buffers Compiler)

### Start the Server

The server listens on **:50055** with a default TTL of **15 seconds**.

```bash
go run server/*.go
```

### Run the Client (Demo)

In a separate terminal, start a demo client:

```bash
go run client/main.go -name=order-service -port=8080
```


## Project Structure

```plaintext
.
├── proto/           # gRPC service definitions (Protobuf)
├── pkg/
│   └── storage/     # Core sharded in-memory engine
├── server/          # gRPC server implementation
└── client/          # Demo client
```


## Future Improvements

* Consistent hashing for multi-node cluster mode
* Watch API using gRPC streaming for real-time service updates
* Persistent storage (Write-Ahead Log) for crash recovery
