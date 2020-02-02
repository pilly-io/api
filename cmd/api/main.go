package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/pilly-io/api/internal/apis"
	"github.com/pilly-io/api/internal/config"
	"github.com/pilly-io/api/internal/models"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	r := gin.New()
	log := logrus.New()
	r.Use(ginlogrus.Logger(log), gin.Recovery())

	config.Settings.DB, config.Settings.DBErr = gorm.Open("sqlite3", ":memory")
	if config.Settings.DBErr != nil {
		panic(config.Settings.DBErr)
	}
	config.Settings.DB.AutoMigrate(&models.Cluster{})

	v1 := r.Group("/api/v1")
	{
		v1.GET("/clusters/:name", apis.GetClusterByName)
		//v1.GET("/clusters/:id", apis.GetClusterById)
	}

	r.Run()
}
