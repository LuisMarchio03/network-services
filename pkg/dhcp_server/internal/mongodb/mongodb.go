package mongodb

import (
	"context"
	"log"
	"network-services-server-dhcp/internal/mongodb/schemas"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB é uma estrutura que contém a conexão com o MongoDB.
type MongoDB struct {
	Client     *mongo.Client     // Client é o cliente MongoDB usado para interagir com o banco de dados.
	Collection *mongo.Collection // Collection é a coleção específica no banco de dados MongoDB.
}

// ConnectMongoDB estabelece uma conexão com o MongoDB utilizando o URI fornecido.
// Parâmetros:
//   - ctx: Contexto de execução para controle de tempo e cancelamento.
//   - uri: URI de conexão com o servidor MongoDB.
//   - dbName: Nome do banco de dados a ser utilizado.
//   - collectionName: Nome da coleção no banco de dados a ser utilizada.
//
// Retorna:
//   - Um ponteiro para a estrutura MongoDB contendo a conexão estabelecida e a coleção especificada.
//   - Um possível erro, se houver problemas durante a conexão.
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

// Close encerra a conexão com o MongoDB.
// Parâmetros:
//   - ctx: Contexto de execução para controle de tempo e cancelamento.
func (db *MongoDB) Close(ctx context.Context) {
	err := db.Client.Disconnect(ctx)
	if err != nil {
		log.Println("Erro ao fechar a conexão com o MongoDB:", err)
	}
}

// InsertIP insere um novo registro de endereço IP no MongoDB.
// Parâmetros:
//   - ctx: Contexto de execução para controle de tempo e cancelamento.
//   - ip: Endereço IP a ser inserido.
//   - assigned: Indica se o endereço IP está atribuído a um cliente.
//
// Retorna:
//   - Um possível erro, se houver problemas durante a inserção.
func (db *MongoDB) InsertIP(ctx context.Context, ip string, assigned bool) error {
	_, err := db.Collection.InsertOne(ctx, bson.M{"ip_address": ip, "assigned": assigned})
	return err
}

// FindAvailableIP busca no MongoDB por um endereço IP disponível que ainda não foi atribuído.
// Retorna o endereço IP encontrado ou um erro se não houver IPs disponíveis.
// Parâmetros:
//   - ctx: Contexto de execução para controle de tempo e cancelamento.
//
// Retorna:
//   - O endereço IP encontrado.
//   - Um possível erro, se houver problemas durante a busca.
func (db *MongoDB) FindAvailableIP(ctx context.Context) (string, error) {
	var ip schemas.IPRecord
	err := db.Collection.FindOne(ctx, bson.M{"assigned": false}).Decode(&ip)
	if err != nil {
		return "", err
	}
	return ip.IPAddress, nil
}

// UpdateIPAssignment atualiza o status de atribuição de um endereço IP no MongoDB.
// Parâmetros:
//   - ctx: Contexto de execução para controle de tempo e cancelamento.
//   - ip: Endereço IP a ser atualizado.
//   - assigned: Novo status de atribuição para o endereço IP.
//
// Retorna:
//   - Um possível erro, se houver problemas durante a atualização.
func (db *MongoDB) UpdateIPAssignment(ctx context.Context, ip string, assigned bool) error {
	_, err := db.Collection.UpdateOne(ctx, bson.M{"ip_address": ip}, bson.M{"$set": bson.M{"assigned": assigned}})
	return err
}

// DeleteIP remove um registro de endereço IP do MongoDB.
// Parâmetros:
//   - ctx: Contexto de execução para controle de tempo e cancelamento.
//   - ip: Endereço IP a ser removido.
//
// Retorna:
//   - Um possível erro, se houver problemas durante a remoção.
func (db *MongoDB) DeleteIP(ctx context.Context, ip string) error {
	_, err := db.Collection.DeleteOne(ctx, bson.M{"ip_address": ip})
	return err
}
