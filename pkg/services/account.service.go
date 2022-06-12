package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/croisade/chimichanga/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type AccountService interface {
	CreateAccount(*models.Account) (*models.Account, error)
	GetAccount(string) (*models.Account, error)
	GetAccounts() ([]*models.Account, error)
	FindByRefreshToken(string) (*models.Account, error)
	DeleteAccount(string) error
	UpdateAccount(*models.Account) (*models.Account, error)
	Login(*LoginValidation) (*models.Account, error)
	Logout(*LogoutValidation) error
}

type AccountServiceImpl struct {
	accountcollection *mongo.Collection
	ctx               context.Context
}

type LoginValidation struct {
	Email    string `json:"email" bson:"email" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required"`
}

type LogoutValidation struct {
	AccountId string `json:"accountId" binding:"required"`
}

func NewAccountServiceImpl(accountcollection *mongo.Collection, ctx context.Context) *AccountServiceImpl {
	return &AccountServiceImpl{
		accountcollection: accountcollection,
		ctx:               ctx,
	}
}

func (s *AccountServiceImpl) CreateAccount(account *models.Account) (*models.Account, error) {
	var result *models.Account

	insertFilter := bson.M{"email": account.Email}
	err := s.accountcollection.FindOne(s.ctx, insertFilter).Decode(&result)
	if err == nil {
		return nil, errors.New("account already exists")
	}

	hashedPassword, err := s.HashPassword(account.Password)
	if err != nil {
		return nil, err
	}
	account.Password = hashedPassword
	account.AccountId = primitive.NewObjectID().Hex()
	account.CreatedAt = primitive.Timestamp{T: uint32(time.Now().Unix())}
	account.UpdatedAt = primitive.Timestamp{T: uint32(time.Now().Unix())}

	_, err = s.accountcollection.InsertOne(s.ctx, account)

	if err != nil {
		return nil, err
	}

	filter := bson.M{"accountId": account.AccountId}
	err = s.accountcollection.FindOne(s.ctx, filter).Decode(&result)
	return result, err
}

func (s *AccountServiceImpl) GetAccount(accountId string) (*models.Account, error) {
	var result *models.Account
	var err error

	filter := bson.M{"accountId": accountId}
	err = s.accountcollection.FindOne(s.ctx, filter).Decode(&result)

	return result, err
}

func (s *AccountServiceImpl) GetAccounts() ([]*models.Account, error) {
	var results []*models.Account
	var err error
	var cursor *mongo.Cursor
	cursor, err = s.accountcollection.Find(s.ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	err = cursor.All(s.ctx, &results)
	return results, err
}

func (s *AccountServiceImpl) DeleteAccount(accountId string) error {
	_, err := s.accountcollection.DeleteOne(s.ctx, bson.M{"accountId": accountId})
	return err
}

func (s *AccountServiceImpl) UpdateAccount(account *models.Account) (*models.Account, error) {
	filter := bson.M{"accountId": account.AccountId}
	var result *models.Account

	existingAccount, err := s.GetAccount(account.AccountId)

	if err != nil {
		return nil, err
	}

	// ? Can users update their email
	// if account.Email != "" {
	// 	existingAccount.Email = account.Email
	// }
	if account.Password != "" {
		existingAccount.Password = account.Password
	}
	if account.FirstName != "" {
		existingAccount.FirstName = account.FirstName
	}
	if account.LastName != "" {
		existingAccount.LastName = account.LastName
	}
	if account.RefreshToken != "" {
		existingAccount.RefreshToken = account.RefreshToken
	}

	existingAccount.UpdatedAt = primitive.Timestamp{T: uint32(time.Now().Unix())}

	// updatedAccount, err := s.accountcollection.UpdateOne(s.ctx, filter, bson.M{"$set": existingAccount})

	upsert := false
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	updatedAccount := s.accountcollection.FindOneAndUpdate(s.ctx, filter, bson.M{"$set": existingAccount}, &opt)
	if updatedAccount.Err() != nil {
		return nil, updatedAccount.Err()
	}

	decodeErr := updatedAccount.Decode(&result)

	return result, decodeErr
}

func (s *AccountServiceImpl) FindByRefreshToken(refreshToken string) (*models.Account, error) {
	var result *models.Account
	var err error

	filter := bson.M{"refreshToken": refreshToken}
	err = s.accountcollection.FindOne(s.ctx, filter).Decode(&result)

	return result, err
}

func (s *AccountServiceImpl) Login(login *LoginValidation) (*models.Account, error) {
	var result *models.Account
	var err error

	hashedPassword, err := s.HashPassword(login.Password)
	if err != nil {
		return nil, err
	}
	fmt.Println(hashedPassword)

	filter := bson.M{"email": login.Email}
	err = s.accountcollection.FindOne(s.ctx, filter).Decode(&result)

	if err != nil {
		return nil, err
	}

	isValid := s.CheckPasswordHash(login.Password, result.Password)

	if isValid != true {
		return nil, errors.New("Invalid password")
	}

	return result, err
}

func (s *AccountServiceImpl) Logout(login *LogoutValidation) error {
	account := &models.Account{AccountId: login.AccountId, RefreshToken: "loggedOut"}
	_, err := s.UpdateAccount(account)

	return err
}

func (s *AccountServiceImpl) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *AccountServiceImpl) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
