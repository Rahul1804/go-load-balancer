package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-load-balancer/internal/handler"
	"go-load-balancer/pkg/log"
)

func StartServer(lb *handler.LoadBalancer, port string) {
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      lb,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error.Fatalf("Could not listen on %s: %v\n", port, err)
		}
	}()
	log.Info.Printf("Server is ready to handle requests at %s", port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
	log.Info.Println("Server stopped")
}
