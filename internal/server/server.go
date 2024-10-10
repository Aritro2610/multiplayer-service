package server

import (
	"context"
	"encoding/json"
	"log"
	pb "multiplayer-service/proto/multiplayer-service/proto"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	pb.UnimplementedMultiplayerModeServiceServer
	db    *mongo.Collection
	cache *redis.Client
}

func NewServer(db *mongo.Collection, cache *redis.Client) *Server {
	return &Server{
		db:    db,
		cache: cache,
	}
}

func (s *Server) GetPopularMode(ctx context.Context, req *pb.MultiplayerModeRequest) (*pb.MultiplayerModeResponse, error) {
	log.Println(("HERERERERER"))
	cachedResponse, err := s.cache.Get(ctx, req.AreaCode).Result()
	if err == nil && cachedResponse != "" {
		var response pb.MultiplayerModeResponse
		// Unmarshal from cache and return
		err = json.Unmarshal([]byte(cachedResponse), &response)
		if err == nil {
			log.Println("Cache hit: returning data from Redis")
			return &response, nil
		}
		return &response, nil
	}
	var result struct {
		ModeName     string `bson:"mode_name"`
		PlayersCount uint32 `bson:"players_count"`
	}
	// err = s.db.FindOne(ctx, map[string]interface{}{"area_code": req.AreaCode}).Decode(&result)
	err = s.db.FindOne(ctx, bson.M{"area_code": req.AreaCode}).Decode(&result)
	if err != nil {
		return nil, err
	}
	// Dummy response
	log.Println("result from route", result)
	response := &pb.MultiplayerModeResponse{
		AreaCode:     req.AreaCode,
		ModeName:     result.ModeName,
		PlayersCount: result.PlayersCount,
	}
	responseJSON, err := json.Marshal(response)
	if err == nil {
		// Set cache with an expiration time
		err = s.cache.Set(ctx, req.AreaCode, responseJSON, 5*time.Minute).Err()
		if err != nil {
			log.Printf("Failed to cache data: %v", err)
		}
	}
	return response, nil
}
