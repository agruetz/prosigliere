// Package main provides the entry point for the server
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/agruetz/prosigliere/internal/datastore/pg"
	"github.com/agruetz/prosigliere/internal/service"
	blogpb "github.com/agruetz/prosigliere/protos/v1/blog"
)

var (
	// gRPC server settings
	grpcPort = flag.Int("grpc-port", 9090, "The gRPC server port")

	// HTTP/REST gateway settings
	httpPort = flag.Int("http-port", 8080, "The HTTP server port")

	// Database settings
	dbHost     = flag.String("db-host", "localhost", "Database host")
	dbPort     = flag.Int("db-port", 5432, "Database port")
	dbUser     = flag.String("db-user", "postgres", "Database user")
	dbPassword = flag.String("db-password", "postgres", "Database password")
	dbName     = flag.String("db-name", "prosigliere", "Database name")
	dbSSLMode  = flag.String("db-sslmode", "disable", "Database SSL mode")
)

func main() {
	flag.Parse()

	// Initialize logger
	logger := log.New(os.Stdout, "", log.LstdFlags)

	// Create a context that can be canceled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize the database connection
	store, err := pg.New(
		pg.WithHost(*dbHost),
		pg.WithPort(*dbPort),
		pg.WithUser(*dbUser),
		pg.WithPassword(*dbPassword),
		pg.WithDatabase(*dbName),
		pg.WithSSLMode(*dbSSLMode),
		pg.WithMaxOpenConns(10),
		pg.WithMaxIdleConns(5),
		pg.WithConnMaxLife(time.Minute*5),
	)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	logger.Println("Connected to database")

	// Create the blog service
	blogService := service.NewBlogService(store)

	// Start the gRPC server
	go startGRPCServer(ctx, logger, blogService)

	// Start the HTTP/REST gateway
	go startHTTPServer(ctx, logger)

	// Wait for termination signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	logger.Println("Received termination signal, shutting down...")
	cancel()

	// Allow some time for graceful shutdown
	time.Sleep(time.Second)
	logger.Println("Server shutdown complete")
}

func startGRPCServer(ctx context.Context, logger *log.Logger, blogService *service.BlogService) {
	addr := fmt.Sprintf(":%d", *grpcPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatalf("Failed to listen on %s: %v", addr, err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the blog service
	blogpb.RegisterBlogsServer(grpcServer, blogService)

	// Register reflection service on gRPC server
	reflection.Register(grpcServer)

	logger.Printf("Starting gRPC server on %s", addr)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Wait for context cancellation to stop the server
	<-ctx.Done()
	logger.Println("Stopping gRPC server...")
	grpcServer.GracefulStop()
	logger.Println("gRPC server stopped")
}

func startHTTPServer(ctx context.Context, logger *log.Logger) {
	addr := fmt.Sprintf(":%d", *httpPort)
	mux := runtime.NewServeMux()

	// Set up a connection to the gRPC server
	grpcAddr := fmt.Sprintf("localhost:%d", *grpcPort)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Register the blog service handler
	err := blogpb.RegisterBlogsHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
	if err != nil {
		logger.Fatalf("Failed to register gateway: %v", err)
	}

	// Create an HTTP server
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	logger.Printf("Starting HTTP/REST gateway on %s", addr)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to serve HTTP: %v", err)
		}
	}()

	// Wait for context cancellation to stop the server
	<-ctx.Done()
	logger.Println("Stopping HTTP server...")

	// Create a deadline for server shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Fatalf("HTTP server shutdown failed: %v", err)
	}
	logger.Println("HTTP server stopped")
}
