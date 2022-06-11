package services

import (
	"context"
	"errors"
	"time"

	"github.com/croisade/chimichanga/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RunService interface {
	CreateRun(*models.Run) (*models.Run, error)
	GetRun(*RunRequest) (*models.Run, error)
	GetAll(*RunFetchRequest) ([]*models.Run, error)
	UpdateRun(*RunUpdateRequest) (*models.Run, error)
	DeleteRun(*RunRequest) error
}

type RunServiceImpl struct {
	runCollection *mongo.Collection
	ctx           context.Context
}

type RunRequest struct {
	AccountId string `json:"accountId" bson:"runId" binding:"required"`
	RunId     string `json:"runId" bson:"runId" binding:"required"`
}

type RunFetchRequest struct {
	AccountId string `json:"accountId" bson:"runId" binding:"required"`
}

type RunUpdateRequest struct {
	AccountId string  `json:"accountId" bson:"runId" binding:"required"`
	RunId     string  `json:"runId" bson:"runId" binding:"required"`
	Pace      float32 `json:"pace" bson:"pace"`
	Time      string  `json:"time" bson:"time"`
	Distance  float32 `json:"distance" bson:"distance"`
	Lap       int     `json:"lap" bson:"lap"`
	Incline   float32 `json:"incline" bson:"incline"`
}

func NewRunService(runCollection *mongo.Collection, ctx context.Context) *RunServiceImpl {
	return &RunServiceImpl{
		runCollection: runCollection,
		ctx:           ctx,
	}
}

func (u *RunServiceImpl) CreateRun(run *models.Run) (*models.Run, error) {
	var result *models.Run

	run.RunId = primitive.NewObjectID().Hex()
	run.CreatedAt = primitive.Timestamp{T: uint32(time.Now().Unix())}
	run.UpdatedAt = primitive.Timestamp{T: uint32(time.Now().Unix())}

	_, err := u.runCollection.InsertOne(u.ctx, result)
	return run, err
}

func (u *RunServiceImpl) GetRun(runRequest *RunRequest) (*models.Run, error) {
	var run *models.Run
	query := bson.M{"accountId": runRequest.AccountId, "runId": runRequest.RunId}

	err := u.runCollection.FindOne(u.ctx, query).Decode(&run)
	return run, err
}

func (u *RunServiceImpl) GetAll(runAccountId *RunFetchRequest) ([]*models.Run, error) {
	var runs []*models.Run
	filter := bson.M{"accountId": runAccountId}

	cursor, err := u.runCollection.Find(u.ctx, filter)

	if err != nil {
		return nil, err
	}

	for cursor.Next(u.ctx) {
		var run models.Run
		err := cursor.Decode(&run)
		if err != nil {
			return nil, err
		}
		runs = append(runs, &run)

	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	cursor.Close(u.ctx)

	if len(runs) == 0 {
		return nil, errors.New("documents not found")
	}
	return runs, nil
}

func (u *RunServiceImpl) UpdateRun(run *RunUpdateRequest) (*models.Run, error) {
	filter := bson.M{"accountId": run.AccountId}
	var result *models.Run

	existingRun, err := u.GetRun(&RunRequest{run.AccountId, run.RunId})

	if err != nil {
		return nil, err
	}

	if run.Time != "" {
		existingRun.Time = run.Time
	}
	if run.Distance != 0.0 {
		existingRun.Distance = run.Distance
	}
	if run.Incline != 0.0 {
		existingRun.Incline = run.Incline
	}
	if run.Lap != 0 {
		existingRun.Lap = run.Lap
	}
	if run.Pace != 0.0 {
		existingRun.Pace = run.Pace
	}

	existingRun.UpdatedAt = primitive.Timestamp{T: uint32(time.Now().Unix())}

	upsert := false
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	updatedRun := u.runCollection.FindOneAndUpdate(u.ctx, filter, bson.M{"$set": existingRun}, &opt)
	if updatedRun.Err() != nil {
		return nil, updatedRun.Err()
	}

	decodeErr := updatedRun.Decode(&result)

	return result, decodeErr
}

func (u *RunServiceImpl) DeleteRun(runRequest *RunRequest) error {
	filter := bson.M{"accountId": runRequest.AccountId, "runId": runRequest.RunId}

	result, _ := u.runCollection.DeleteOne(u.ctx, filter)

	if result.DeletedCount != 1 {
		return errors.New("no matched document found for delete")
	}

	return nil
}
