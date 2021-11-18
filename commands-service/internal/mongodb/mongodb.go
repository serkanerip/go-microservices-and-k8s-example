package mongodb

import (
	"context"
	"log"

	"github.com/serkanerip/commands-service/config"
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
	Id         primitive.ObjectID `bson:"_id,omitempty"`
	ExternalId string             `bson:"external_id"`
	Name       string             `bson:"name"`
}

type Command struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	HowTo       string             `bson:"how_to"`
	CommandLine string             `bson:"command_line"`
	PlatformId  primitive.ObjectID `bson:"platform_id"`
}
