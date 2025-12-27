# Go High-Performance Service Registry

A lightweight, high-concurrency Service Registry & Discovery server built with **Go**, **gRPC**, and **Sharded Architecture**. Designed to handle millions of operations per second with minimal lock contention.

## Key Features

- **High Performance Storage**: Implements a **Sharded Map (32 shards)** to reduce `sync.RWMutex` contention significantly compared to standard maps.
- **gRPC & Protobuf**: Uses efficient binary serialization for low-latency communication.
- **Service Governance**: Supports Server-Side TTL (Time-To-Live) management and active heartbeat monitoring.
- **Thread-Safe**: Fully concurrent-safe design optimized for mixed read-write workloads.

## Architecture

The core engine uses a **Sharding Strategy** (similar to ConcurrentHashMap in Java) to split the storage into multiple buckets based on FNV-1a hash of the service name.

```mermaid
graph TD
    Client -->|Register/Discover| Server[gRPC Server]
    Server -->|Hash(Key)| ShardSelector
    ShardSelector -->|Index=0| Shard0[Shard 0 (Lock)]
    ShardSelector -->|Index=1| Shard1[Shard 1 (Lock)]
    ShardSelector -->|...| ShardN[Shard ... (Lock)]
    ShardSelector -->|Index=31| Shard31[Shard 31 (Lock)]
```

## Benchmark Results

TBD

## Quick Start
**Prerequisites**
* Go 1.21+
* Protoc (Protocol Buffers Compiler)

**1. Start Server**
The server will start on port :50055 with a default TTL of 15 seconds.

```bash
go run server/*.go
```

**2. Run Client (Demo)**
In a new terminal, run the client simulation.

```bash
go run client/main.go -name=order-service -port=8080
```

## Project Structure
```plaintext
.
├── proto/           # gRPC Service Definitions (Protobuf)
├── pkg/
│   └── storage/     # Core Sharded Map Engine (The "Brain")
├── server/          # gRPC Server Implementation
└── client/          # Demo Client SDK
```

## Future Improvements
- Implement Consistent Hashing for cluster mode.
- Add Watch API (gRPC Streaming) for real-time updates.
- Persistent storage (Write-Ahead Log) for crash recovery.
