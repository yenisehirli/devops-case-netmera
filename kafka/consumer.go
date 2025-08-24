package kafka

import (
    "context"
    "encoding/json"
    "log"
    "time"
    
    "go-kafka-app/models"
    "go-kafka-app/mongo"
    
    "github.com/IBM/sarama"
)

type Consumer struct {
    consumer sarama.Consumer
    topic    string
    mongoClient *mongo.Client
}

func NewConsumer(brokers []string, topic string, mongoClient *mongo.Client) (*Consumer, error) {
    config := sarama.NewConfig()
    config.Consumer.Return.Errors = true
    
    consumer, err := sarama.NewConsumer(brokers, config)
    if err != nil {
        return nil, err
    }
    
    return &Consumer{
        consumer:    consumer,
        topic:       topic,
        mongoClient: mongoClient,
    }, nil
}

func (c *Consumer) Start(ctx context.Context) error {
    partitions, err := c.consumer.Partitions(c.topic)
    if err != nil {
        return err
    }
    
    for _, partition := range partitions {
        partitionConsumer, err := c.consumer.ConsumePartition(c.topic, partition, sarama.OffsetNewest)
        if err != nil {
            log.Printf("Error creating partition consumer: %v", err)
            continue
        }
        
        go func(pc sarama.PartitionConsumer) {
            defer pc.Close()
            
            for {
                select {
                case message := <-pc.Messages():
                    c.processMessage(message)
                case err := <-pc.Errors():
                    log.Printf("Kafka consumer error: %v", err)
                case <-ctx.Done():
                    return
                }
            }
        }(partitionConsumer)
    }
    
    return nil
}

func (c *Consumer) processMessage(msg *sarama.ConsumerMessage) {
    log.Printf("Received message from Kafka, waiting 10 seconds...")
    
    // 10 saniye bekle
    time.Sleep(10 * time.Second)
    
    var messageDoc models.MessageDocument
    if err := json.Unmarshal(msg.Value, &messageDoc); err != nil {
        log.Printf("Error unmarshaling message: %v", err)
        return
    }
    
    // MongoDB'ye kaydet
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := c.mongoClient.InsertMessage(ctx, messageDoc); err != nil {
        log.Printf("Error saving to MongoDB: %v", err)
        return
    }
    
    log.Printf("Message processed successfully: %s", messageDoc.ID)
}

func (c *Consumer) Close() error {
    return c.consumer.Close()
}
