package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/carldunham/useful-cookery/internal/ai"
	"github.com/carldunham/useful-cookery/internal/auth"
	"github.com/carldunham/useful-cookery/internal/config"
	"github.com/carldunham/useful-cookery/internal/database"
	"github.com/carldunham/useful-cookery/internal/graphql"
	"github.com/carldunham/useful-cookery/internal/graphql/resolvers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := log.New(os.Stdout, "", log.LstdFlags)
	logger.Printf("Starting Useful Cookery API server on port %s", cfg.Server.Port)

	// Connect to DGraph
	dgraphClient, err := database.NewDGraphClient(cfg.DGraph.Hosts)
	if err != nil {
		logger.Fatalf("Failed to connect to DGraph: %v", err)
	}
	logger.Println("Connected to DGraph")

	// Initialize cache
	cache := database.NewRedisCache(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)

	// Initialize AI service
	aiConfig := &ai.Config{
		OpenAIAPIKey:     cfg.AI.OpenAIKey,
		EmbeddingModel:   cfg.AI.EmbeddingModel,
		CompletionModel:  cfg.AI.CompletionModel,
		CacheEnabled:     cfg.AI.EnableCache,
		CacheTTL:         time.Duration(cfg.AI.CacheTTLMinutes) * time.Minute,
		MaxRequestTokens: cfg.AI.MaxTokens,
	}

	aiService, err := ai.NewAIService(aiConfig, cache)
	if err != nil {
		logger.Fatalf("Failed to initialize AI service: %v", err)
	}
	logger.Println("AI service initialized")

	// Initialize auth service
	authService := auth.NewService(
		cfg.Auth.JWTSecret,
		cfg.Auth.TokenExpiry,
		dgraphClient,
	)

	// Initialize GraphQL resolvers
	resolver := resolvers.NewRootResolver(dgraphClient, aiService, authService)

	// Create GraphQL handler
	gqlHandler := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{
		Resolvers: resolver,
	}))

	// Create router
	router := chi.NewRouter()

	// Add middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	// Add CORS middleware
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.Server.CorsOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Define routes
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Useful Cookery API"))
	})

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// GraphQL endpoint
	router.Post("/graphql", authService.AuthMiddleware(gqlHandler.ServeHTTP))

	// GraphQL playground
	if cfg.Server.EnablePlayground {
		router.Get("/playground", playground.Handler("GraphQL playground", "/graphql"))
		logger.Println("GraphQL Playground available at /playground")
	}

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeoutSeconds) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeoutSeconds) * time.Second,
	}

	// Start server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server error: %v", err)
		}
	}()

	logger.Printf("Server listening on port %s", cfg.Server.Port)

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("Shutting down server...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Println("Server exited gracefully")
}
