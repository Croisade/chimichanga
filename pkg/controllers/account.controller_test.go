package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

func TestLogin(t *testing.T) {
	want := &models.Account{Email: "test@example.com", Password: "password", FirstName: "first", LastName: "last"}
	token := &JWTtoken{}
	accountService.CreateAccount(want)

	r := SetupRouter()
	r.GET("/account/login", accountController.Login)

	login := &services.Login{Email: "test@example.com", Password: "password"}
	jsonValue, _ := json.Marshal(login)
	req, _ := http.NewRequest("GET", "/account/login", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), token)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 3, len(strings.Split(token.Token, ".")))
}
