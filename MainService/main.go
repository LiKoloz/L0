package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	. "L0_WB/controllers"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	server := &http.Server{
		Addr:    ":8081",
		Handler: nil, // Используем DefaultServeMux
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		select {
		case <-time.After(30 * time.Second):
			APP_PORT, exists := os.LookupEnv("APP_PORT")
			if !exists {
				log.Println("APP_PORT not set in environment")
				return
			}

			req, err := http.NewRequestWithContext(
				ctx,
				"GET",
				"http://localhost:"+APP_PORT+"/order/order_888_it",
				nil,
			)
			if err != nil {
				log.Printf("Failed to create request: %v", err)
				return
			}

			client := &http.Client{Timeout: 10 * time.Second}
			resp, err := client.Do(req)
			if err != nil {
				if ctx.Err() != nil {
					log.Println("Request cancelled due to shutdown")
				} else {
					log.Printf("Request error: %v", err)
				}
				return
			}
			defer resp.Body.Close()

			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Error reading response: %v", err)
				return
			}
			log.Printf("Server response: %s", string(bodyBytes))

		case <-ctx.Done():
			log.Println("Scheduled request cancelled due to shutdown")
			return
		}
	}()

	OrderEndPoint()

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Server starting on :8081")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v\n", err)
		}
	}()

	sig := <-quit
	log.Printf("Received signal: %v. Starting graceful shutdown...\n", sig)

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v\n", err)
	} else {
		log.Println("HTTP server stopped gracefully")
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("All background tasks completed")
	case <-time.After(5 * time.Second):
		log.Println("Timeout waiting for background tasks, forcing exit")
	}

	log.Println("Graceful shutdown completed")
}
