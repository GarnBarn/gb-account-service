package main

import (
	"fmt"
	"github.com/GarnBarn/common-go/database"
	"github.com/GarnBarn/common-go/httpserver"
	"github.com/GarnBarn/common-go/logger"
	"github.com/GarnBarn/gb-account-service/config"
	"github.com/gin-contrib/cors"
	"github.com/sirupsen/logrus"
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
	_, err := database.Conn(appConfig.MYSQL_CONNECTION_STRING)
	if err != nil {
		logrus.Panic("Can't connect to db: ", err)
	}

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

	logrus.Info("Listening and serving HTTP on :", appConfig.HTTP_SERVER_PORT)
	httpServer.Run(fmt.Sprint(":", appConfig.HTTP_SERVER_PORT))
}
