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
	_, err := db.Collection.InsertOne(ctx, bson.M{"ip_address": ip, "assigned": assigned})
	return err
}

// FindAvailableIP busca no MongoDB por um endereço IP disponível que ainda não foi atribuído.
// Retorna o endereço IP encontrado ou um erro se não houver IPs disponíveis.
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

// ReleaseIPAddress marca um endereço IP como disponível novamente no MongoDB
func (db *MongoDB) ReleaseIPAddress(ctx context.Context, ip string) error {
	// Atualiza o status de atribuição do endereço IP para 'false' (não atribuído) no MongoDB
	err := db.UpdateIPAssignment(ctx, ip, false)
	if err != nil {
		return err
	}

	// Registro de log indicando que o endereço IP foi liberado com sucesso
	log.Printf("Endereço IP %s foi liberado com sucesso", ip)

	return nil
}
