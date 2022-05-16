package services

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/croisade/chimichanga/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	want := &models.User{}
	createdAt := primitive.Timestamp{
		T: uint32(time.Now().Unix()),
	}
	usr := &models.User{UserId: "user123", CreatedAt: createdAt, UpdatedAt: createdAt, Runs: models.Runs{Pace: 5.0, Time: "5:00", Distance: "5 miles", Date: createdAt, Lap: 2, SessionId: "123", UserId: "user123"}}

	userService := NewUserService(col, ctx)
	got, err := userService.CreateUser(usr)
	assert.Nil(t, err)

	query := bson.D{bson.E{Key: "userId", Value: "user123"}}
	findErr := col.FindOne(ctx, query).Decode(&want)

	assert.Nil(t, findErr)
	assert.Equal(t, want, got)
}
