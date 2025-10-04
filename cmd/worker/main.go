package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ecommerce_project/internal/config"
	"ecommerce_project/pkg/db"
	"ecommerce_project/pkg/logger"
)

func main() {
	// Initialize logger
	logger.Init()
	defer logger.Sync()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	database, err := db.NewConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start background workers
	go startEmailWorker(ctx)
	go startNotificationWorker(ctx)
	go startInventoryWorker(ctx)
	go startOrderProcessingWorker(ctx)

	logger.Info("Background workers started")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down workers...")
	cancel()

	// Give workers time to cleanup
	time.Sleep(5 * time.Second)
	logger.Info("Workers stopped")
}

func startEmailWorker(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Email worker stopped")
			return
		case <-ticker.C:
			// Process email queue
			logger.Debug("Processing email queue...")
		}
	}
}

func startNotificationWorker(ctx context.Context) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Notification worker stopped")
			return
		case <-ticker.C:
			// Process notification queue
			logger.Debug("Processing notification queue...")
		}
	}
}

func startInventoryWorker(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Inventory worker stopped")
			return
		case <-ticker.C:
			// Check and update inventory
			logger.Debug("Checking inventory levels...")
		}
	}
}

func startOrderProcessingWorker(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Order processing worker stopped")
			return
		case <-ticker.C:
			// Process pending orders
			logger.Debug("Processing pending orders...")
		}
	}
}
