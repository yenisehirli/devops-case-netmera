package mongo

import (
    "context"
    "log"
    "time"
    
    "go-kafka-app/models"
    
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
    client     *mongo.Client
    database   *mongo.Database
    collection *mongo.Collection
}

func NewClient(uri, database, collection string) (*Client, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
        return nil, err
    }
    
    db := client.Database(database)
    coll := db.Collection(collection)
    
    return &Client{
        client:     client,
        database:   db,
        collection: coll,
    }, nil
}

func (c *Client) InsertMessage(ctx context.Context, message models.MessageDocument) error {
    message.ProcessedAt = time.Now()
    
    _, err := c.collection.InsertOne(ctx, message)
    if err != nil {
        log.Printf("Error inserting message: %v", err)
        return err
    }
    
    log.Printf("Message saved to MongoDB: %s", message.ID)
    return nil
}

func (c *Client) Close() error {
    return c.client.Disconnect(context.Background())
}
