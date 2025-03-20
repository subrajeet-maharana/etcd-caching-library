# Cache System

A modular and scalable cache system with multiple implementations and `etcd` integration.

## Features

- Multiple cache implementations:
  - BTree-based cache for efficient range queries
  - LRU cache for automatic eviction based on usage patterns
- `etcd` integration for distributed cache synchronization
- Benchmarking utilities for performance evaluation
- Logging system with rotation capabilities
- Comprehensive test suite

## Architecture

The system is organized into the following components:

- `package/cache`: Contains cache implementations (BTree, LRU) and the common interface
- `package/etcd`: `etcd` client for distributed cache synchronization
- `package/logging`: Logging utilities with rotation support
- `package/benchmark`: Benchmarking tools for cache performance evaluation
- `cmd/main.go`: Main application entry point

## Getting Started

### Prerequisites

- Go 1.23+
- `etcd` cluster (for distributed features)

### Installation

```bash
# Clone the repository
git clone https://github.com/subrajeet-maharana/etcd-caching-library.git
cd etcd-caching-library

# Install dependencies
make deps

# Build the application
make build

# Run tests
make test

# Run the application
make run
```

### Docker

```bash
# Build Docker image
make docker-build

# Run Docker container
make docker-run
```

## Usage

```go
import (
    "etcd-caching-library/package/cache"
    "etcd-caching-library/package/etcd"
)

func main() {
    // Initialize cache
    btreeCache := cache.NewBTreeCache()

    // Set and get values
    btreeCache.Set("key", "value")
    value, found := btreeCache.Get("key")

    // List values in a range
    items := btreeCache.List("start", "end")

    // Use LRU cache
    lruCache := cache.NewLRUCache(100)
    lruCache.Set("key", "value")
    value, found = lruCache.Get("key")
}
```

## Benchmarking

The package includes benchmarking tools to evaluate cache performance:

```go
import "etcd-caching-library/package/benchmark"

// Benchmark BTree cache
benchmark.BenchmarkBTreeCache(btreeCache, numReaders, numWriters, iterations)

// Benchmark LRU cache
benchmark.BenchmarkLRUCache(lruCache, numReaders, numWriters, iterations)
```
