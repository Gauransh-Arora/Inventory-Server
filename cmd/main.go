package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"server/internal/config"
	"server/internal/handler"
	"server/internal/repository"
	"server/internal/service"
	"server/internal/utils"
	"server/routes"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: Error loading .env file (using system environment variables)")
	}

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	config.ConnectDB()
	defer config.DB.Close()

	if err := utils.InitJWT(); err != nil {
		log.Fatal("Failed to initialize JWT keys: ", err)
	}

	repoInv := repository.NewInventoryRepository(config.DB)
	repoProd := repository.NewProductRepository(config.DB)
	repoAuth := repository.NewAuthRepository(config.DB)

	serviceInv := service.NewInventoryService(repoInv)
	serviceProd := service.NewProductService(repoProd)
	serviceAuth := service.NewAuthService(repoAuth)

	handlerInv := handler.NewInventoryHandler(serviceInv)
	handlerProd := handler.NewProductHandler(serviceProd)
	handlerAuth := handler.NewAuthHandler(serviceAuth)

	service.StartCleanupTask(repoAuth)

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	routes.SetupRoutes(r, handlerInv, handlerProd, handlerAuth, repoAuth)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("Server listening on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exited cleanly")
}