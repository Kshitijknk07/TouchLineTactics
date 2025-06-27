package storage

import (
	"context"
	"os"

	"github.com/yourusername/TouchlineTactics/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func GetMongoClient() (*mongo.Client, error) {
	if mongoClient != nil {
		return mongoClient, nil
	}
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	mongoClient = client
	return client, nil
}

func FetchRandomPlayers(n int) ([]domain.Player, error) {
	client, err := GetMongoClient()
	if err != nil {
		return nil, err
	}
	coll := client.Database("auction").Collection("players")
	pipeline := mongo.Pipeline{
		{{Key: "$sample", Value: bson.D{{Key: "size", Value: n}}}},
	}
	cur, err := coll.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	var players []domain.Player
	for cur.Next(context.Background()) {
		var p domain.Player
		if err := cur.Decode(&p); err == nil {
			players = append(players, p)
		}
	}
	return players, nil
}

func FetchRandomPlayersByPosition(position string, n int) ([]domain.Player, error) {
	client, err := GetMongoClient()
	if err != nil {
		return nil, err
	}
	coll := client.Database("auction").Collection("players")
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "Position", Value: position}}}},
		{{Key: "$sample", Value: bson.D{{Key: "size", Value: n}}}},
	}
	cur, err := coll.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	var players []domain.Player
	for cur.Next(context.Background()) {
		var p domain.Player
		if err := cur.Decode(&p); err == nil {
			players = append(players, p)
		}
	}
	return players, nil
}
