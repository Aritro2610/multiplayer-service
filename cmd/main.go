// main.go
package main

import (
	"context"
	"log"
	"multiplayer-service/internal/server"
	pb "multiplayer-service/proto/multiplayer-service/proto"
	"net"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

// type server struct {
// 	pb.UnimplementedMultiplayerModeServiceServer
// 	db    *mongo.Collection
// 	cache *redis.Client
// }

// func (s *server) GetPopularMode(ctx context.Context, req *pb.MultiplayerModeRequest) (*pb.MultiplayerModeResponse, error) {
// 	cachedResponse, err := s.cache.Get(ctx, req.AreaCode).Result()
// 	if err == nil && cachedResponse != "" {
// 		var response pb.MultiplayerModeResponse
// 		// Unmarshal from cache and return
// 		return &response, nil
// 	}
// 	var result struct {
// 		ModeName     string `bson:"mode_name"`
// 		PlayersCount uint32 `bson:"players_count"`
// 	}
// 	err = s.db.FindOne(ctx, map[string]interface{}{"area_code": req.AreaCode}).Decode(&result)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// Dummy response
// 	response := &pb.MultiplayerModeResponse{
// 		AreaCode:     "123",
// 		ModeName:     "Battle Royale",
// 		PlayersCount: 150,
// 		// AreaCode:     req.AreaCode,
// 		// ModeName:     result.ModeName,
// 		// PlayersCount: result.PlayersCount,
// 	}
// 	s.cache.Set(ctx, req.AreaCode, response, 5*time.Minute)
// 	return response, nil
// }

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Connect to MongoDB
	// log.Println(os.Getenv("MONGO_URI"))
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")

	db := client.Database("multiplayerDB").Collection("modes")
	// mode := models.MultiplayerMode{
	//     ModeName:   "Battle Royale",
	//     AreaCode:   "123",
	//     PlayerCount: 150,
	//     LastUpdated: time.Now().Unix(),
	// }

	// _, err = db.InsertOne(ctx, mode)
	// if err != nil {
	//     log.Fatalf("Failed to insert document: %v", err)
	// }
	// log.Println("Document inserted, collection and database will be created if they don't exist.")

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis address
		DB:   0,                // use default DB
	})

	s := server.NewServer(db, redisClient)

	// Register gRPC server
	pb.RegisterMultiplayerModeServiceServer(grpcServer, s)

	log.Println("Server is running on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
