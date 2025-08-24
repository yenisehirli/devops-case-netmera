package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "strings"
    "syscall"
    
    "go-kafka-app/handlers"
    "go-kafka-app/kafka"
    "go-kafka-app/mongo"
    
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

func main() {
    fmt.Println("Starting Go Kafka App...")
    
    // Environment variables yükle
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system environment variables")
    }
    
    // Configuration
    kafkaBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
    kafkaTopic := os.Getenv("KAFKA_TOPIC")
    mongoURI := os.Getenv("MONGODB_URI")
    mongoDatabase := os.Getenv("MONGODB_DATABASE")
    
    // MongoDB client başlat
    mongoClient, err := mongo.NewClient(mongoURI, mongoDatabase, "messages")
    if err != nil {
        log.Fatalf("Failed to create MongoDB client: %v", err)
    }
    defer mongoClient.Close()
    
    // Kafka producer başlat
    producer, err := kafka.NewProducer(kafkaBrokers, kafkaTopic)
    if err != nil {
        log.Fatalf("Failed to create Kafka producer: %v", err)
    }
    defer producer.Close()
    
    // Kafka consumer başlat
    consumer, err := kafka.NewConsumer(kafkaBrokers, kafkaTopic, mongoClient)
    if err != nil {
        log.Fatalf("Failed to create Kafka consumer: %v", err)
    }
    defer consumer.Close()
    
    // Consumer'ı background'da başlat
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    if err := consumer.Start(ctx); err != nil {
        log.Fatalf("Failed to start Kafka consumer: %v", err)
    }
    
    // Handler'a producer'ı ver
    handlers.InitKafkaProducer(producer)
    
    // HTTP server başlat
    router := gin.Default()
    
    router.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })
    
    router.POST("/message", handlers.PostMessage)
    
    // Graceful shutdown için
    go func() {
        log.Println("Server starting on :8080")
        if err := router.Run(":8080"); err != nil {
            log.Fatalf("Server failed to start: %v", err)
        }
    }()
    
    // SIGINT/SIGTERM bekle
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    log.Println("Shutting down server...")
    cancel()
}
