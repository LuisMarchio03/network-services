package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB é a estrutura que contém a conexão com o banco de dados MongoDB.
type MongoDB struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

// NewMongoDB ConnectMongoDB, é a func responsavel por estabelecer a conexão com  o mongodb utilizando um URI
//
// Como parâmentros essa func vai receber:
// -- ctx: Contexto de execução para controle de tempo e cancelamento do mongoDB
// -- uri: URI necessaria para realizar a conexão com o banco de dados (mongoDB)
// -- dbName: Representa o nome do banco de dados em execução dentro do servidor
// -- collectionName: Nome da collection no banco de dados que está a ser utilizado
//
// Como retorno essa func vai devolver:
// -- Um ponteiro para a estrutura MongoDB contendo a conexão já estabelecida e a collection especificada
// -- Um possível erro, se houver problemas durante a tentativa de conexão com o banco de dados
func NewMongoDB(ctx context.Context, uri, dbName, collectionName string) (*MongoDB, error) {
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
