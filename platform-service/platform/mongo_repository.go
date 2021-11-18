package platform

import (
	"context"
	"log"

	"github.com/serkanerip/platform-service/internal/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

var ctx = context.TODO()

func NewMongoDBRepository(client *mongo.Client) *MongoDBRepository {
	return &MongoDBRepository{
		client:     client,
		collection: client.Database("db").Collection("platforms"),
	}
}

func (m *MongoDBRepository) GetAll() ([]PlatformDTO, error) {
	dtos := make([]PlatformDTO, 0)
	cur, err := m.collection.Find(ctx, bson.D{{}})
	if err != nil {
		log.Printf("cannot get platforms err is: %v", err)
		return nil, err
	}

	for cur.Next(ctx) {
		var p mongodb.Platform
		if err := cur.Decode(&p); err != nil {
			return nil, err
		}
		dtos = append(dtos, *PersistenceToPlatformDTO(p))
	}

	return dtos, nil
}

func (m *MongoDBRepository) GetPlatformById(id string) (*PlatformDTO, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": bson.M{"$eq": objID}}

	var p mongodb.Platform
	if err := m.collection.FindOne(ctx, filter).Decode(&p); err != nil {
		return nil, err
	}

	return PersistenceToPlatformDTO(p), nil
}

func (m *MongoDBRepository) CreatePlatform(dto CreatePlatformDTO) (string, error) {

	insertResult, err := m.collection.InsertOne(ctx, mongodb.Platform{
		Name:      dto.Name,
		Publisher: dto.Publisher,
		Cost:      dto.Cost,
	})
	if err != nil {
		return "", err
	}

	return insertResult.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (m *MongoDBRepository) Seed() {
	count, err := m.collection.CountDocuments(ctx, bson.D{{}})
	if err != nil {
		log.Printf("Cannot get document count err is: %v", err)
	}
	if count > 0 {
		log.Println("--> We already have data")
		return
	}

	log.Println("--> Seeding Data...")
	m.collection.InsertMany(ctx, []interface{}{
		mongodb.Platform{Name: "Go", Publisher: "Google", Cost: "Free"},
		mongodb.Platform{Name: "Mysql", Publisher: "Mysql", Cost: "Free"},
		mongodb.Platform{Name: "Kubernetes", Publisher: "Cloud Native Computing Foundation", Cost: "Free"},
	})
}
