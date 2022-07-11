package services

import (
	"context"
	"errors"
	"strconv"
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
	AccountId string `json:"accountId" bson:"accountId" binding:"required"`
	RunId     string `json:"runId" bson:"runId" binding:"required"`
}

type RunFetchRequest struct {
	AccountId string `json:"accountId" bson:"accountId" binding:"required"`
	Date      string `json:"date" bson:"date"`
}

type RunUpdateRequest struct {
	AccountId string  `json:"accountId" bson:"accountId" binding:"required"`
	RunId     string  `json:"runId" bson:"runId" binding:"required"`
	Speed     float32 `json:"speed" bson:"speed"`
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
	run.CreatedAt = time.Now()
	run.UpdatedAt = time.Now()

	_, err := u.runCollection.InsertOne(u.ctx, run)

	if err != nil {
		return nil, err
	}

	filter := bson.M{"accountId": run.AccountId}
	err = u.runCollection.FindOne(u.ctx, filter).Decode(&result)
	return result, err
}

func (u *RunServiceImpl) GetRun(runRequest *RunRequest) (*models.Run, error) {
	var run *models.Run
	query := bson.M{"accountId": runRequest.AccountId, "runId": runRequest.RunId}

	err := u.runCollection.FindOne(u.ctx, query).Decode(&run)
	return run, err
}

func (u *RunServiceImpl) GetAll(runFetchRequest *RunFetchRequest) ([]*models.Run, error) {
	var runs []*models.Run
	var filter primitive.M
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"createdAt", -1}})

	if runFetchRequest.Date != "" {
		i, err := strconv.ParseInt(runFetchRequest.Date, 10, 64)
		if err != nil {
			panic(err)
		}
		filter = bson.M{"accountId": runFetchRequest.AccountId, "createdAt": bson.M{
			"$gt": time.Unix(i/1000, 0),
			"$lt": time.Unix((i/1000)+86400, 0),
		}}

	} else {
		filter = bson.M{"accountId": runFetchRequest.AccountId}
	}

	cursor, err := u.runCollection.Find(u.ctx, filter, findOptions)

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
	if run.Speed != 0.0 {
		existingRun.Speed = run.Speed
	}

	existingRun.UpdatedAt = time.Now()

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
