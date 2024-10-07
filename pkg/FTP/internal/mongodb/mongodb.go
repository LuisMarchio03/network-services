package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
)

var mongoDBInstance *MongoDB
var once sync.Once

type MongoDB struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func GetMongoDBInstance() (*MongoDB, error) {
	var err error
	once.Do(func() {
		clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
		client, err := mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			log.Fatalf("Erro ao conectar ao MongoDB: %v", err)
		}
		mongoDBInstance = &MongoDB{
			Client:     client,
			Collection: client.Database("ftpdb").Collection("users"),
		}
	})
	return mongoDBInstance, err
}

func (db *MongoDB) FindUser(ctx context.Context, username, password string) (bool, error) {
	var result bson.M
	err := db.Collection.FindOne(ctx, bson.M{"username": username, "password": password}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("erro ao procurar usuário no MongoDB: %v", err)
	}
	return true, nil
}
