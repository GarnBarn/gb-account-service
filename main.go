package main

import (
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	"github.com/GarnBarn/common-go/database"
	"github.com/GarnBarn/common-go/httpserver"
	"github.com/GarnBarn/common-go/logger"
	"github.com/GarnBarn/gb-account-service/config"
	"github.com/GarnBarn/gb-account-service/handler"
	"github.com/GarnBarn/gb-account-service/repository"
	"github.com/GarnBarn/gb-account-service/service"
	"github.com/gin-contrib/cors"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"time"
)

var appConfig config.Config

func init() {
	appConfig = config.Load()
	logger.InitLogger(logger.Config{
		Env: appConfig.Env,
	})

}

func main() {
	// Start DB Connection
	db, err := database.Conn(appConfig.MYSQL_CONNECTION_STRING)
	if err != nil {
		logrus.Panic("Can't connect to db: ", err)
	}

	// Initilize the Firebase App
	opt := option.WithCredentialsFile(appConfig.FIREBASE_CONFIG_FILE)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logrus.Fatalln("error initializing app: %v\n", err)
	}

	// Init repository
	accountRepository := repository.NewAccountRepository(db)

	// Init service
	accountService := service.NewAccountService(app, accountRepository)

	// Init handler
	accountHandler := handler.NewAccountHandler(accountService)

	// Create the http server
	httpServer := httpserver.NewHttpServer()

	httpServer.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	router := httpServer.Group("/api/v1")
	accountRouter := router.Group("/account")
	accountRouter.GET("/", accountHandler.GetAccount)

	logrus.Info("Listening and serving HTTP on :", appConfig.HTTP_SERVER_PORT)
	httpServer.Run(fmt.Sprint(":", appConfig.HTTP_SERVER_PORT))
}
