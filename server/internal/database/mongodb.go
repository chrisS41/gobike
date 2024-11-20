package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	COL_NAME_RIDES  = "rides"
	COL_NAME_ROUTES = "routes"
	COL_NAME_USERS  = "users"
)

type Collection struct {
	collection *mongo.Collection
}

type MongoDB struct {
	client *mongo.Client
	db     *mongo.Database
	Rides  *Collection
	Routes *Collection
	Users  *Collection
}

func NewMongoDB(uri, dbName string) (*MongoDB, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)

	return &MongoDB{
		client: client,
		db:     db,
		Rides:  &Collection{collection: db.Collection(COL_NAME_RIDES)},
		Routes: &Collection{collection: db.Collection(COL_NAME_ROUTES)},
		Users:  &Collection{collection: db.Collection(COL_NAME_USERS)},
	}, nil
}

// Collection 구조체의 메서드들
func (c *Collection) Create(document interface{}) (primitive.ObjectID, error) {
	result, err := c.collection.InsertOne(context.Background(), document)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func (c *Collection) ReadOne(filter interface{}, result interface{}) error {
	err := c.collection.FindOne(context.Background(), filter).Decode(result)
	if err != nil {
		return err
	}
	return nil
}

func (c *Collection) ReadMany(filter interface{}) ([]bson.M, error) {
	cursor, err := c.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (c *Collection) Update(filter interface{}, update interface{}) error {
	_, err := c.collection.UpdateOne(
		context.Background(),
		filter,
		update,
	)
	return err
}

func (c *Collection) Delete(filter interface{}) error {
	_, err := c.collection.DeleteOne(context.Background(), filter)
	return err
}

// MongoDB 연결 종료
func (m *MongoDB) Close() error {
	return m.client.Disconnect(context.Background())
}
