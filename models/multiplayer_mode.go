package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type MultiplayerMode struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	ModeName    string             `bson:"mode_name"`
	AreaCode    string             `bson:"area_code"`
	PlayerCount int                `bson:"players_count"`
	LastUpdated int64              `bson:"last_updated"`
}
