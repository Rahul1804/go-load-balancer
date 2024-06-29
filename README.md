# Go Load Balancer

This project implements a simple load balancer in Go, which uses round-robin load balancing and includes health checks for backend servers. It uses Go's `http` package for handling HTTP requests and the `sync` package for concurrency control. 

## Features

- Round-robin load balancing
- Health checks for backend servers
- Concurrency-safe implementation
- Logging of requests and errors
- Graceful shutdown
- Panic recovery
