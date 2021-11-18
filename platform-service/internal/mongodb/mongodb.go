package mongodb

import (
	"context"
	"log"

	"github.com/serkanerip/platform-service/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.TODO()

func CreateClient() *mongo.Client {
	clientOptions := options.Client().ApplyURI(config.ENV.MONGODB_CONNECTION_STRING)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("cannot connect mongodb err is: %v", err)
	}

	return client
}

type Platform struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Publisher string             `bson:"publisher"`
	Cost      string             `bson:"cost"`
}
