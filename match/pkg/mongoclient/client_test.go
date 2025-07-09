package mongoclient_test

import (
	"context"
	"match/pkg/mongoclient"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name"`
	Email string             `bson:"email"`
	Age   int                `bson:"age"`
}

func getTestConfig() *mongoclient.Config {
	return &mongoclient.Config{
		Hosts:        []string{"localhost:27017"},
		DatabaseName: "testdb",
		MinPoolSize:  1,
		MaxPoolSize:  5,
	}
}

func TestMongoClient_CRUD(t *testing.T) {
	conf := getTestConfig()
	client := mongoclient.New(conf)

	err := client.Start()
	if err != nil {
		t.Fatalf("Failed to start client: %v", err)
	}
	defer client.Close()

	collection := client.Collection("users")
	ctx := context.Background()

	// --- Create ---
	user := User{
		Name:  "Alice",
		Email: "alice@example.com",
		Age:   30,
	}

	insertRes, err := collection.InsertOne(ctx, user)
	if err != nil {
		t.Fatalf("InsertOne failed: %v", err)
	}
	insertedID := insertRes.InsertedID.(primitive.ObjectID)

	// --- Read ---
	var result User
	err = collection.FindOne(ctx, bson.M{"_id": insertedID}).Decode(&result)
	if err != nil {
		t.Fatalf("FindOne failed: %v", err)
	}
	if result.Name != user.Name {
		t.Errorf("Expected name %s, got %s", user.Name, result.Name)
	}

	// --- Update ---
	update := bson.M{
		"$set": bson.M{"age": 35},
	}
	_, err = collection.UpdateOne(ctx, bson.M{"_id": insertedID}, update)
	if err != nil {
		t.Fatalf("UpdateOne failed: %v", err)
	}

	// Confirm Update
	var updated User
	err = collection.FindOne(ctx, bson.M{"_id": insertedID}).Decode(&updated)
	if err != nil {
		t.Fatalf("FindOne after update failed: %v", err)
	}
	if updated.Age != 35 {
		t.Errorf("Expected age 35, got %d", updated.Age)
	}

	// --- Delete ---
	_, err = collection.DeleteOne(ctx, bson.M{"_id": insertedID})
	if err != nil {
		t.Fatalf("DeleteOne failed: %v", err)
	}

	// Confirm Delete
	err = collection.FindOne(ctx, bson.M{"_id": insertedID}).Decode(&User{})
	if err == nil {
		t.Errorf("Expected error on FindOne after delete, got nil")
	}
}
