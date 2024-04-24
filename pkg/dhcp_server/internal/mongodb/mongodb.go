package mongodb

import (
	"context"
	"log"
	"network-services-server-dhcp/internal/mongodb/schemas"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB é uma estrutura que contém a conexão com o MongoDB
type MongoDB struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

// ConnectMongoDB estabelece uma conexão com o MongoDB
func ConnectMongoDB(ctx context.Context, uri, dbName, collectionName string) (*MongoDB, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	collection := client.Database(dbName).Collection(collectionName)

	return &MongoDB{
		Client:     client,
		Collection: collection,
	}, nil
}

// Close encerra a conexão com o MongoDB
func (db *MongoDB) Close(ctx context.Context) {
	err := db.Client.Disconnect(ctx)
	if err != nil {
		log.Println("Erro ao fechar a conexão com o MongoDB:", err)
	}
}

// InsertIP insere um novo registro de endereço IP no MongoDB
func (db *MongoDB) InsertIP(ctx context.Context, ip string, assigned bool) error {
	_, err := db.Collection.InsertOne(ctx, schemas.IPRecord{IPAddress: ip, Assigned: assigned})
	return err
}

// FindAvailableIP busca no MongoDB por um endereço IP disponível que ainda não foi atribuído.
// Funcionamento Geral:
// A função procura no banco de dados MongoDB por um documento que tenha o campo assigned com valor false,
// indicando que o endereço IP correspondente não foi atribuído a nenhum cliente.
//
// Parâmetros:
// ctx context.Context: O contexto da execução, usado para controlar o tempo de execução e o cancelamento de operações.
//
// Retorno:
// string: O endereço IP encontrado.
// error: Qualquer erro que ocorra durante a execução da consulta.
func (db *MongoDB) FindAvailableIP(ctx context.Context) (string, error) {
	var ip schemas.IPRecord
	err := db.Collection.FindOne(ctx, bson.M{"assigned": false}).Decode(&ip)
	if err != nil {
		return "", err
	}
	return ip.IPAddress, nil
}

// UpdateIPAssignment atualiza o status de atribuição de um endereço IP no MongoDB
func (db *MongoDB) UpdateIPAssignment(ctx context.Context, ip string, assigned bool) error {
	_, err := db.Collection.UpdateOne(ctx, bson.M{"ip_address": ip}, bson.M{"$set": bson.M{"assigned": assigned}})
	return err
}

// DeleteIP remove um registro de endereço IP do MongoDB
func (db *MongoDB) DeleteIP(ctx context.Context, ip string) error {
	_, err := db.Collection.DeleteOne(ctx, bson.M{"ip_address": ip})
	return err
}
