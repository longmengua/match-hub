package mongoclient

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID    string `bson:"_id,omitempty"`
	Name  string `bson:"name"`
	Email string `bson:"email"`
	Age   int    `bson:"age"`
}

type MongoClient struct {
	client     *mongo.Client
	collection *mongo.Collection
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewMongoClient(uri, dbName, collectionName string) (*MongoClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		cancel()
		return nil, err
	}

	return &MongoClient{
		client:     client,
		collection: client.Database(dbName).Collection(collectionName),
		ctx:        ctx,
		cancel:     cancel,
	}, nil
}

func (m *MongoClient) Start() error {
	return m.client.Ping(m.ctx, nil)
}

func (m *MongoClient) Close() {
	m.cancel()
	_ = m.client.Disconnect(m.ctx)
}

func (m *MongoClient) CreateUser(user User) (*mongo.InsertOneResult, error) {
	return m.collection.InsertOne(m.ctx, user)
}

func (m *MongoClient) GetUserByName(name string) (*User, error) {
	var result User
	err := m.collection.FindOne(m.ctx, bson.M{"name": name}).Decode(&result)
	return &result, err
}

func (m *MongoClient) UpdateUserAge(name string, age int) (*mongo.UpdateResult, error) {
	return m.collection.UpdateOne(m.ctx, bson.M{"name": name}, bson.M{"$set": bson.M{"age": age}})
}

func (m *MongoClient) DeleteUser(name string) (*mongo.DeleteResult, error) {
	return m.collection.DeleteOne(m.ctx, bson.M{"name": name})
}
