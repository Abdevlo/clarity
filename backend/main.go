package main

import (
	"fmt"
	"log"
	"net"

	"github.com/clarity/backend/config"
	"github.com/clarity/backend/database"
	authpb "github.com/clarity/backend/gen/go/auth"
	healthpb "github.com/clarity/backend/gen/go/health"
	aipb "github.com/clarity/backend/gen/go/ai"
	"github.com/clarity/backend/handlers"
	"github.com/clarity/backend/services"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	log.Printf("Starting server on %s:%s", cfg.Server.Host, cfg.Server.Port)

	// Initialize database
	db, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := db.Migrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	dbConn := db.GetConnection()
	defer db.Close()

	// Initialize services
	authService := services.NewAuthService(dbConn, &cfg.Auth)
	healthService := services.NewHealthRecordsService(dbConn)
	aiService := services.NewAIService(dbConn)

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register services
	authpb.RegisterAuthServiceServer(grpcServer, handlers.NewAuthServer(authService))
	healthpb.RegisterHealthRecordsServiceServer(grpcServer, handlers.NewHealthRecordsServer(healthService))
	aipb.RegisterAIServiceServer(grpcServer, handlers.NewAIServer(aiService))

	// Listen on port
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("gRPC server listening on %s:%s", cfg.Server.Host, cfg.Server.Port)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
