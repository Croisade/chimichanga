package services

import (
	"context"
	"log"
	"testing"

	"github.com/croisade/chimichanga/pkg/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var accountCollection *mongo.Collection
var runsCollection *mongo.Collection
var ctx context.Context

func setup() {
	ctx = context.TODO()
	c, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	accountCollection = c.Database("CorroYouRun").Collection("accounts")
	runsCollection = c.Database("CorroYouRun").Collection("runs")
	_, delErr := runsCollection.DeleteMany(ctx, bson.D{{}})
	_, accDelErr := accountCollection.DeleteMany(ctx, bson.D{{}})

	if delErr != nil {
		log.Fatalf("Error: %v", delErr)
	}
	if accDelErr != nil {
		log.Fatalf("Error: %v", delErr)
	}
	log.Println("\n-----Setup complete-----")
}

func TestMain(m *testing.M) {
	setup()
	m.Run()
	// teardown()
}

func TestCreateRun(t *testing.T) {
	run := &models.Run{Speed: 6.0, Lap: 0, Distance: 3.0, Time: "30:00", Incline: 0.0, AccountId: "123"}
	response := &models.Run{}

	runService := NewRunService(runsCollection, ctx)
	got, err := runService.CreateRun(run)
	assert.Nil(t, err)

	query := bson.M{"accountId": run.AccountId}
	findErr := runsCollection.FindOne(ctx, query).Decode(&response)

	assert.Nil(t, findErr)
	assert.Equal(t, response.Speed, got.Speed)
}
