package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/croisade/chimichanga/pkg/models"
	"github.com/croisade/chimichanga/pkg/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	gin.SetMode(gin.TestMode)
	// {
	// 	userGroup := v1.Group("user")
	// 	{
	// 		account := NewAccountController()
	// 		userGroup.GET("/:id", user.Retrieve)
	// 	}
	// }
	return router
}

var accountCollection *mongo.Collection
var accountController AccountController
var accountService *services.AccountServiceImpl
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

	accountService = services.NewAccountServiceImpl(accountCollection, ctx)
	jwtService := services.NewJWTAuthService()

	accountController = NewAccountController(accountService, jwtService)
	log.Println("\n-----Setup complete-----")
}

func TestMain(m *testing.M) {
	setup()
	m.Run()
	// teardown()
}

type Response struct {
	Errors []ErrorValidationMsg `json:"errors"`
}

func TestCreate(t *testing.T) {
	t.Run("Should create an account", func(t *testing.T) {
		account := &models.Account{Email: "test@example.com", Password: "password", FirstName: "first", LastName: "last"}
		response := &models.Account{}
		r := SetupRouter()
		r.POST("/account/create", accountController.CreateAccount)

		jsonValue, _ := json.Marshal(account)
		req, _ := http.NewRequest("POST", "/account/create", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		json.Unmarshal(w.Body.Bytes(), response)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, account.Email, response.Email)
	})

	t.Run("Should raise if a required field is missing", func(t *testing.T) {
		account := &models.Account{Email: "test@example.com", Password: "password", FirstName: "first"}
		response := &Response{}

		r := SetupRouter()
		r.POST("/account/create", accountController.CreateAccount)

		jsonValue, _ := json.Marshal(account)
		req, _ := http.NewRequest("POST", "/account/create", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, response.Errors[0].Field, "LastName")
		assert.Equal(t, response.Errors[0].Message, "This field is required")
	})
}

func TestLogin(t *testing.T) {
	t.Run("Should Raise if account does not exist", func(t *testing.T) {
		accountCollection.DeleteMany(ctx, bson.D{{}})
		response := &ErrorResponse{}

		r := SetupRouter()
		r.PUT("/account/login", accountController.Login)

		login := &services.LoginValidation{Email: "test@example.com", Password: "password"}
		jsonValue, _ := json.Marshal(login)
		req, _ := http.NewRequest("PUT", "/account/login", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		json.Unmarshal(w.Body.Bytes(), response)
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, "mongo: no documents in result", response.Errors)
	})

	t.Run("Should create a token", func(t *testing.T) {
		want := &models.Account{Email: "test@example.com", Password: "password", FirstName: "first", LastName: "last"}
		accountService.CreateAccount(want)
		token := &JWTtoken{}

		r := SetupRouter()
		r.PUT("/account/login", accountController.Login)

		login := &services.LoginValidation{Email: "test@example.com", Password: "password"}
		jsonValue, _ := json.Marshal(login)
		req, _ := http.NewRequest("PUT", "/account/login", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		json.Unmarshal(w.Body.Bytes(), token)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, 3, len(strings.Split(token.Token, ".")))
	})

	t.Run("Should validate a request", func(t *testing.T) {
		login := &services.LoginValidation{Email: "test@example.com"}

		response := &Response{}

		r := SetupRouter()
		r.PUT("/account/login", accountController.Login)

		jsonValue, _ := json.Marshal(login)
		req, _ := http.NewRequest("PUT", "/account/login", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, response.Errors[0].Field, "Password")
		assert.Equal(t, response.Errors[0].Message, "This field is required")
	})
}

func TestToken(t *testing.T) {
	t.Run("Should refresh a token", func(t *testing.T) {
		var refreshToken JWTtoken
		var response *JWTtoken
		token, _ := accountController.JWTService.CreateToken()

		want := &models.Account{Email: "test@example.com", Password: "password", FirstName: "first", LastName: "last", RefreshToken: token}
		accountService.CreateAccount(want)
		refreshToken.RefreshToken = token

		time.Sleep(1 * time.Second)

		r := SetupRouter()
		r.PUT("/account/token", accountController.Token)

		jsonValue, _ := json.Marshal(refreshToken)
		req, _ := http.NewRequest("PUT", "/account/token", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotEqual(t, response.Token, token)
	})

	t.Run("Should raise if a token is missing", func(t *testing.T) {
		response := &ErrorResponse{}

		r := SetupRouter()
		r.PUT("/account/token", accountController.Token)

		req, _ := http.NewRequest("PUT", "/account/token", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, response.Errors, "invalid request")
	})
}
func TestLogout(t *testing.T) {
	t.Run("Should logout and remove token from account document", func(t *testing.T) {
		want := &models.Account{Email: "test@example.com", Password: "password", FirstName: "first", LastName: "last"}
		accountService.CreateAccount(want)
		accounts, _ := accountService.GetAccounts()
		var logout = &services.LogoutValidation{AccountId: accounts[0].AccountId}

		// token, _ := accountController.JWTService.CreateToken()

		// accountCollection.FindOne(ctx, {})

		r := SetupRouter()
		r.PUT("/account/logout", accountController.Logout)

		jsonValue, _ := json.Marshal(logout)
		req, _ := http.NewRequest("PUT", "/account/logout", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestGetAccount(t *testing.T) {
	var response *models.Account
	fixture := &models.Account{Email: "test@example.com", Password: "password", FirstName: "first", LastName: "last"}
	account, _ := accountService.CreateAccount(fixture)

	r := SetupRouter()
	r.GET("/account/get/:accountId", accountController.GetAccount)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/account/get/%v", account.AccountId), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, response.AccountId, account.AccountId)
}

func TestGetAccounts(t *testing.T) {
	accountCollection.DeleteMany(ctx, bson.D{{}})
	var response []*models.Account
	fixture := &models.Account{Email: "test@example.com", Password: "password", FirstName: "first", LastName: "last"}
	account, _ := accountService.CreateAccount(fixture)

	r := SetupRouter()
	r.GET("/account/fetch", accountController.GetAccounts)

	req, _ := http.NewRequest("GET", "/account/fetch", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, response[0].AccountId, account.AccountId)
}
func TestDeleteAccount(t *testing.T) {
	fixture := &models.Account{Email: "test@example.com", Password: "password", FirstName: "first", LastName: "last"}
	account, _ := accountService.CreateAccount(fixture)

	r := SetupRouter()
	r.DELETE("/account/delete/:accountId", accountController.DeleteAccount)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/account/delete/%v", account.AccountId), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateAccount(t *testing.T) {
	fixture := &models.Account{Email: "test@example.com", Password: "password", FirstName: "first", LastName: "last"}
	account, _ := accountService.CreateAccount(fixture)
	update := &models.Account{AccountId: account.AccountId, FirstName: "Jamal"}

	var response *models.Account
	jsonValue, _ := json.Marshal(update)

	r := SetupRouter()
	r.PUT("/account/update", accountController.UpdateAccount)

	req, _ := http.NewRequest("PUT", "/account/update", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, response.FirstName, update.FirstName)
	assert.Equal(t, response.LastName, fixture.LastName)
}
