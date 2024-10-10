package test

import (
	"context"
	"log"
	"multiplayer-service/internal/server"
	"testing"
	"time"

	pb "multiplayer-service/proto/multiplayer-service/proto"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func TestGetPopularMode(t *testing.T) {

	// Connect to Redis
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis address for local testing
	})
	_, errRedis := redisClient.Ping(context.Background()).Result()
	if errRedis != nil {
		t.Fatalf("Failed to connect to Redis: %v", errRedis)
	}
	defer redisClient.FlushAll(context.Background()) // Clear Redis after the test

	// Connect to MongoDB
	// mongoURI := "mongodb://root:example@localhost:27017" // MongoDB URI for local testing
	clientOptions := options.Client().ApplyURI("mongodb+srv://yash:yashMongodb@ecommerce-backend.eyunzek.mongodb.net/?retryWrites=true&w=majority&appName=Ecommerce-Backend")

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		t.Fatalf("Failed to create MongoDB client: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Ping MongoDB to verify connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		t.Fatalf("Failed to ping MongoDB: %v", err)
	}
	log.Println("Connected to MongoDB")

	// Use the appropriate database and collection
	db := client.Database("multiplayerDB").Collection("modes")

	// Ensure the collection is cleaned up before/after the test
	// defer func() {
	// 	if err := db.Drop(ctx); err != nil {
	// 		t.Fatalf("Failed to drop test collection: %v", err)
	// 	}
	// }()

	// // Insert a test document into MongoDB
	// testData := bson.M{
	// 	"area_code":     "123",
	// 	"mode_name":     "Battle Royale",
	// 	"players_count": 150,
	// }
	// _, err = db.InsertOne(ctx, testData)
	// if err != nil {
	// 	t.Fatalf("Failed to insert test document: %v", err)
	// }

	// Create a new server with the actual db and cache
	s := server.NewServer(db, redisClient)

	// Create a request
	req := &pb.MultiplayerModeRequest{AreaCode: "14"}

	// Call GetPopularMode
	res, err := s.GetPopularMode(context.Background(), req)
	log.Println(res)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check response
	if res.ModeName != "Battle Royale" {
		t.Errorf("Expected 'Battle Royale', got %v", res.ModeName)
	}

	if res.PlayersCount != 150 {
		t.Errorf("Expected '150', got %v", res.PlayersCount)
	}

	log.Println("Test passed.", req)
}
