package schemas

import "go.mongodb.org/mongo-driver/bson/primitive"

// IPRecord representa um registro de endereço IP no MongoDB.
type IPRecord struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"` // ID é o identificador único do registro no MongoDB.
	IPAddress string             `bson:"ip_address"`    // IPAddress é o endereço IP associado ao registro.
	Assigned  bool               `bson:"assigned"`      // Assigned indica se o endereço IP está atribuído a um cliente.
}
