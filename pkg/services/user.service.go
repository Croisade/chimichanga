package services

import (
	"context"
	"errors"

	"github.com/croisade/chimichanga/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService interface {
	CreateUser(*models.User) (*models.User, error)
	GetUser(*string) (*models.User, error)
	GetAll() ([]*models.User, error)
	UpdateUser(*models.User) error
	DeleteUser(*string) error
}

type UserServiceImpl struct {
	usercollection *mongo.Collection
	ctx            context.Context
}

func NewUserService(usercollection *mongo.Collection, ctx context.Context) *UserServiceImpl {
	return &UserServiceImpl{
		usercollection: usercollection,
		ctx:            ctx,
	}
}

func (u *UserServiceImpl) CreateUser(user *models.User) (*models.User, error) {
	_, err := u.usercollection.InsertOne(u.ctx, user)
	return user, err
}

func (u *UserServiceImpl) GetUser(userId *string) (*models.User, error) {
	var user *models.User
	query := bson.D{bson.E{Key: "userId", Value: userId}}

	err := u.usercollection.FindOne(u.ctx, query).Decode(&user)
	return user, err
}

func (u *UserServiceImpl) GetAll() ([]*models.User, error) {
	var users []*models.User

	cursor, err := u.usercollection.Find(u.ctx, bson.D{{}})

	if err != nil {
		return nil, err
	}

	for cursor.Next(u.ctx) {
		var user models.User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)

	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	cursor.Close(u.ctx)

	if len(users) == 0 {
		return nil, errors.New("documents not found")
	}
	return users, nil
}

func (u *UserServiceImpl) UpdateUser(user *models.User) error {
	filter := bson.D{bson.E{Key: "userId", Value: user.UserId}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{bson.E{Key: "userId", Value: user.UserId}, bson.E{Key: "runs", Value: user.Runs}}}}
	result, _ := u.usercollection.UpdateOne(u.ctx, filter, update)

	if result.MatchedCount != 1 {
		return errors.New("no matched document found for update")
	}
	return nil
}

func (u *UserServiceImpl) DeleteUser(userId *string) error {
	filter := bson.D{bson.E{Key: "userId", Value: userId}}
	result, _ := u.usercollection.DeleteOne(u.ctx, filter)

	if result.DeletedCount != 1 {
		return errors.New("no matched document found for delete")
	}

	return nil
}
