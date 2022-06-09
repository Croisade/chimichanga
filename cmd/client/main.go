package main

import (
	"context"
	"fmt"
	"log"

	"github.com/croisade/chimichanga/pkg/conf"
	"github.com/croisade/chimichanga/pkg/controllers"
	"github.com/croisade/chimichanga/pkg/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server      *gin.Engine
	ctx         context.Context
	mongoClient *mongo.Client
	err         error

	jwtservice services.JWTService

	accountService    services.AccountService
	accountController controllers.AccountController
	accountCollection *mongo.Collection

	usercollection *mongo.Collection
	userservice    services.UserService
	usercontroller controllers.UserController
)

func init() {
	config, err := conf.LoadConfig("../../")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	//? Can I plug the the config into the context?
	ctx = context.TODO()

	mongoConn := options.Client().ApplyURI(config.MongoURI)
	mongoClient, err = mongo.Connect(ctx, mongoConn)

	if err != nil {
		log.Fatal(err)
	}

	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("mongo connection established")

	accountCollection = mongoClient.Database("CorroYouRun").Collection("accounts")
	usercollection = mongoClient.Database("CorroYouRun").Collection("runs")

	jwtService := services.NewJWTAuthService()

	userservice = services.NewUserService(usercollection, ctx)
	usercontroller = controllers.New(userservice)

	accountService = services.NewAccountServiceImpl(accountCollection, ctx)
	accountController = controllers.NewAccountController(accountService, jwtService)

	server = gin.Default()
}

func main() {
	defer mongoClient.Disconnect(ctx)

	basePath := server.Group("/v1")
	usercontroller.RegisterUserRoutes(basePath)
	accountController.RegisterAccountRoutes(basePath)
	log.Fatal(server.Run(":9090"))
}
