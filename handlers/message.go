package handlers

import (
	"fmt"
	"net/http"
	"time"

	"go-kafka-app/kafka"
	"go-kafka-app/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var kafkaProducer *kafka.Producer

func InitKafkaProducer(producer *kafka.Producer) {
	kafkaProducer = producer
}

func PostMessage(c *gin.Context) {
	var request models.MessageRequest

	// JSON'u parse et
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid JSON format",
			"details": err.Error(),
		})
		return
	}

	// Unique ID oluştur
	messageID := uuid.New().String()

	// Kafka için message hazırla
	kafkaMessage := models.MessageDocument{
		ID:        messageID,
		Content:   request.Content,
		UserID:    request.UserID,
		CreatedAt: time.Now(),
	}

	// Kafka'ya gönder
	if err := kafkaProducer.SendMessage(messageID, kafkaMessage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to send message to Kafka",
			"details": err.Error(),
		})
		return
	}

	fmt.Printf("Message sent to Kafka: %+v\n", kafkaMessage)

	// Response hazırla
	response := models.MessageResponse{
		ID:      messageID,
		Content: request.Content,
		UserID:  request.UserID,
		Status:  "sent_to_kafka",
		Time:    time.Now(),
	}

	c.JSON(http.StatusOK, response)
}
