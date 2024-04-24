package schemas

import "go.mongodb.org/mongo-driver/bson/primitive"

// IPRecord representa um registro de endereço IP no MongoDB
type IPRecord struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	IPAddress string             `bson:"ip_address"`
	Assigned  bool               `bson:"assigned"`
}
