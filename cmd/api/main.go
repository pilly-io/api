package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pilly-io/api/internal/apis"
	idb "github.com/pilly-io/api/internal/db"
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

	database, err := idb.New("sqlite3", "/tmp/pilly.sqlite")
	if err != nil {
		panic(err)
	}
	database.Migrate()
	apis.SetupRouter(r, database)
	r.Run()
}
