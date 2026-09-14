package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kontursvet-api/internal/config"
	"kontursvet-api/internal/handler"
	"kontursvet-api/internal/repository"
	"kontursvet-api/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "kontursvet-api/docs"
)

// @title КонтурСвет API Swagger

// @version 1.0
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

func main() {
	cfg := config.Load()

	pool, err := service.GetDBPool(cfg)
	if err != nil {
		log.Fatalf("db pool error: %v", err)
	}
	defer pool.Close()

	leadRepo := repository.NewLeadRepository(pool)
	portfolioRepo := repository.NewPortfolioRepository(pool)

	leadService := service.NewLeadService(leadRepo, cfg)
	maxBotService := service.NewMaxBotService(
		cfg.MaxBot.Token,
		cfg.MaxBot.ChatID,
		cfg.MaxBot.Active,
	)

	leadHandler := handler.NewLeadHandler(leadService, maxBotService)
	portfolioHandler := handler.NewPortfolioHandler(portfolioRepo)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.POST("/api/leads", leadHandler.Create)
	router.GET("/api/leads", leadHandler.List)
	router.POST("/api/upload", leadHandler.Upload)
	router.POST("/api/upload-json", leadHandler.UploadJSON)

	router.GET("/api/portfolio", portfolioHandler.ListCases)
	router.GET("/api/portfolio/:id", portfolioHandler.GetCase)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	go func() {
		fmt.Printf("Server started on port %s\n", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	fmt.Println("Server stopped")
}
