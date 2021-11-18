package command

import (
	"context"

	"github.com/serkanerip/commands-service/internal/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommandRepo interface {
	GetAllPlatforms() ([]Platform, error)
	CreatePlatform(Platform) error
	PlatformExists(string) (bool, error)
	ExternalPlatformExists(string) (bool, error)

	GetCommandsForPlatform(string) ([]Command, error)
	GetCommand(platformId string, commandId string) (*Command, error)
	CreateCommand(string, Command) error
}

type CommandMongoRepository struct {
	client              *mongo.Client
	commandsCollection  *mongo.Collection
	platformsCollection *mongo.Collection
}

var ctx = context.TODO()

func NewCommandMongoRepository(client *mongo.Client) *CommandMongoRepository {
	return &CommandMongoRepository{
		client:              client,
		commandsCollection:  client.Database("db").Collection("commands"),
		platformsCollection: client.Database("db").Collection("platforms"),
	}
}

func (c *CommandMongoRepository) GetAllPlatforms() ([]Platform, error) {
	cur, err := c.platformsCollection.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, nil
	}
	platforms := make([]Platform, 0)
	for cur.Next(ctx) {
		var platform mongodb.Platform
		if err := cur.Decode(&platform); err != nil {
			return nil, nil
		}
		platforms = append(platforms, Platform{
			Id:         platform.Id.Hex(),
			ExternalId: platform.ExternalId,
			Name:       platform.Name,
		})
	}

	return platforms, nil
}

func (c *CommandMongoRepository) CreatePlatform(platform Platform) error {
	_, err := c.platformsCollection.InsertOne(ctx, mongodb.Platform{
		ExternalId: platform.ExternalId,
		Name:       platform.Name,
	})
	return err
}

func (c *CommandMongoRepository) PlatformExists(id string) (bool, error) {
	pid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}
	filter := bson.M{"_id": bson.M{"$eq": pid}}
	count, err := c.platformsCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (c *CommandMongoRepository) ExternalPlatformExists(id string) (bool, error) {
	filter := bson.M{"_id": bson.M{"$eq": id}}
	count, err := c.platformsCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (c *CommandMongoRepository) GetCommandsForPlatform(platformId string) ([]Command, error) {
	pid, err := primitive.ObjectIDFromHex(platformId)
	if err != nil {
		return nil, err
	}
	cur, err := c.commandsCollection.Find(ctx, bson.M{"platform_id": pid})
	if err != nil {
		return nil, err
	}
	commands := make([]Command, 0)
	for cur.Next(ctx) {
		var command mongodb.Command
		if err = cur.Decode(&command); err != nil {
			return nil, err
		}
		commands = append(commands, Command{
			Id:          command.Id.Hex(),
			HowTo:       command.HowTo,
			CommandLine: command.CommandLine,
			PlatformId:  command.PlatformId.Hex(),
		})
	}

	return commands, nil
}

func (c *CommandMongoRepository) GetCommand(platformId string, commandId string) (*Command, error) {
	pid, err := primitive.ObjectIDFromHex(platformId)
	if err != nil {
		return nil, err
	}
	cid, err := primitive.ObjectIDFromHex(commandId)
	if err != nil {
		return nil, err
	}
	res := c.commandsCollection.FindOne(ctx, bson.M{"platform_id": pid, "id": cid})

	if res.Err() != nil {
		return nil, err
	}

	var command Command
	if err := res.Decode(&command); err != nil {
		return nil, err
	}

	return &command, nil
}

func (c *CommandMongoRepository) CreateCommand(platformId string, command Command) error {
	oid, err := primitive.ObjectIDFromHex(platformId)
	if err != nil {
		return err
	}
	_, err = c.commandsCollection.InsertOne(ctx, mongodb.Command{
		HowTo:       command.HowTo,
		CommandLine: command.CommandLine,
		PlatformId:  oid,
	})
	return err
}
