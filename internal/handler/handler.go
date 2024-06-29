package handler

import (
	"io"
	"net/http"
	"sync"
	"time"

	"go-load-balancer/pkg/log"
)

type LoadBalancer struct {
	servers       []string
	index         uint64
	serverStatus  map[string]bool // Map to store server health status
	serverTimeout time.Duration   // Timeout for health check requests
	mu            sync.RWMutex    // Mutex for serverStatus map access
}

func NewLoadBalancer(servers []string, healthCheckTimeout time.Duration) *LoadBalancer {
	lb := &LoadBalancer{
		servers:       servers,
		serverStatus:  make(map[string]bool),
		serverTimeout: healthCheckTimeout,
	}

	// Initialize serverStatus map with healthy status
	for _, server := range servers {
		lb.serverStatus[server] = true
	}

	// Start a background health checker routine
	go lb.healthChecker()

	return lb
}

func (lb *LoadBalancer) RoundRobin() string {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	if len(lb.servers) == 0 {
		return ""
	}
	i := lb.index % uint64(len(lb.servers))
	lb.index++
	return lb.servers[i]
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Error.Printf("Recovered from panic: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}()

	start := time.Now()
	server := lb.RoundRobin()

	if server == "" {
		http.Error(w, "No backend servers available", http.StatusServiceUnavailable)
		return
	}

	lb.mu.RLock()
	healthy := lb.serverStatus[server]
	lb.mu.RUnlock()

	if !healthy {
		http.Error(w, "Backend server is not healthy", http.StatusServiceUnavailable)
		return
	}

	proxyReq, err := http.NewRequest(r.Method, server+r.URL.Path, r.Body)
	if err != nil {
		log.Error.Printf("Error creating request: %v", err)
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	proxyReq.Header = r.Header
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		log.Error.Printf("Error forwarding request: %v", err)
		http.Error(w, "Error forwarding request", http.StatusInternalServerError)
		return
	}
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Error.Printf("Error copying response body: %v", err)
	}

	log.Info.Printf("Forwarded request to %s in %v", server, time.Since(start))
}

func (lb *LoadBalancer) healthChecker() {
	ticker := time.NewTicker(5 * time.Second) // Adjust the interval based on your needs

	for range ticker.C {
		lb.mu.Lock()
		for _, server := range lb.servers {
			if !lb.isServerHealthy(server) {
				lb.serverStatus[server] = false
			} else {
				lb.serverStatus[server] = true
			}
		}
		lb.mu.Unlock()
	}
}

func (lb *LoadBalancer) isServerHealthy(server string) bool {
	client := &http.Client{
		Timeout: lb.serverTimeout,
	}

	req, err := http.NewRequest("GET", server+"/health", nil)
	if err != nil {
		log.Error.Printf("Error creating health check request: %v", err)
		return false
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error.Printf("Error performing health check request: %v", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error.Printf("Server %s returned non-OK status: %v", server, resp.Status)
		return false
	}

	return true
}
