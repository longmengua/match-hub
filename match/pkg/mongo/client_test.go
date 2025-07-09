package mongo_test

import (
	"match/pkg/mongo"
	"testing"
)

func TestMongoClient_CRUD(t *testing.T) {
	client, err := mongo.NewMongoClient("mongodb://localhost:27017", "testdb", "users")
	if err != nil {
		t.Fatalf("連接失敗: %v", err)
	}
	defer client.Close()

	if err := client.Start(); err != nil {
		t.Fatalf("啟動失敗: %v", err)
	}

	// Create
	user := mongo.User{Name: "Bob", Email: "bob@example.com", Age: 20}
	_, err = client.CreateUser(user)
	if err != nil {
		t.Fatalf("新增失敗: %v", err)
	}

	// Read
	got, err := client.GetUserByName("Bob")
	if err != nil || got.Name != "Bob" {
		t.Fatalf("查詢失敗: %v", err)
	}

	// Update
	_, err = client.UpdateUserAge("Bob", 21)
	if err != nil {
		t.Fatalf("更新失敗: %v", err)
	}

	got, _ = client.GetUserByName("Bob")
	if got.Age != 21 {
		t.Fatalf("更新年齡錯誤，應為 21，實為 %d", got.Age)
	}

	// Delete
	_, err = client.DeleteUser("Bob")
	if err != nil {
		t.Fatalf("刪除失敗: %v", err)
	}

	// Ensure deleted
	_, err = client.GetUserByName("Bob")
	if err == nil {
		t.Fatalf("應該找不到資料，但仍查到")
	}
}
