package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/hoangtm1601/go-binance-crawler/internal/api/repositories"
	"github.com/hoangtm1601/go-binance-crawler/internal/api/services"

	"golang.org/x/net/context"

	"github.com/gin-gonic/gin"
	_ "github.com/hoangtm1601/go-binance-crawler/docs"
	"github.com/hoangtm1601/go-binance-crawler/internal/initializers"
	"github.com/hoangtm1601/go-binance-crawler/internal/middleware"
)

var (
	server *gin.Engine
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)
	initializers.InitRedis(&config)
	if config.EnableAutoMigrate == "true" {
		if err := initializers.Migrate(); err != nil {
			log.Fatal("Failed to run database migrations", err)
		}
	}

	server = gin.Default()
}

func main() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	// add timestamp to name to avoid overwrite this log
	f, _ := os.Create("tmp/gin.log")
	gin.DefaultWriter = io.MultiWriter(f)
	server.Use(middleware.Recover())
	server.NoRoute(func(c *gin.Context) {
		c.HTML(404, "404.html", gin.H{})
	})

	// To implementing graceful shutdown
	go func() {
		// Initialize logger
		logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

		// Initialize Binance client
		binanceClient := binance.NewClient(config.BinanceAPIKey, config.BinanceAPISecret)
		// Initialize CandleRepository
		candleRepo := repositories.NewCandleRepository(initializers.DB)
		// Initialize CandleService (assuming you have a constructor for it)
		candleService := services.NewCandleService(candleRepo)

		// Initialize CrawlersService
		crawlersService := services.NewCrawlersService(binanceClient, logger, candleService, &config)

		// Run the Crawl function
		crawlersService.Crawl()
	}()

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
	<-quit
	log.Println("Shutting down crawler...")

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer func() {
		log.Println("Close all connections")
		cancel()
	}()

	close(quit)
}
