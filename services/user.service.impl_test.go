package services

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/croisade/chimichanga/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var col *mongo.Collection
var ctx context.Context

func setup() {
	ctx = context.TODO()
	c, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	col = c.Database("userdb").Collection("users")
	_, delErr := col.DeleteMany(ctx, bson.D{{}})

	if delErr != nil {
		log.Fatalf("Error: %v", delErr)
	}
	log.Println("\n-----Setup complete-----")
}

func TestMain(m *testing.M) {
	setup()
	m.Run()
	// teardown()
}
func TestCreateUser(t *testing.T) {
	userService := NewUserService(col, ctx)
	userService.CreateUser(&models.User{UserId: "user123", Runs: models.Runs{Pace: 5.0, Time: "5:00", Distance: "5 miles", Date: time.Now(), Lap: 2, SessionId: "123", UserId: "user123"}})
	got, _ := col.CountDocuments(ctx, bson.D{})

	if got != 1 {
		t.Errorf("got %q, want 1", got)
	}
}
