package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	accountService    services.AccountService
	accountController controllers.AccountController
	accountCollection *mongo.Collection

	runCollection *mongo.Collection
	runService    services.RunService
	runController controllers.RunController
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
	runCollection = mongoClient.Database("CorroYouRun").Collection("runs")

	jwtService := services.NewJWTAuthService()

	runService = services.NewRunService(runCollection, ctx)
	runController = controllers.NewRunController(runService)

	accountService = services.NewAccountServiceImpl(accountCollection, ctx)
	accountController = controllers.NewAccountController(accountService, jwtService)

	server = gin.Default()
}

func main() {
	defer mongoClient.Disconnect(ctx)

	basePath := server.Group("/v1")
	runController.RegisterRunRoutes(basePath)
	accountController.RegisterAccountRoutes(basePath)

	srv := &http.Server{
		Addr:    ":9090",
		Handler: server,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
