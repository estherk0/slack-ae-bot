package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/estherk0/slack-ae-bot/pkg/db"
	"github.com/estherk0/slack-ae-bot/pkg/routes"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer func() {
		db.Close(ctx)
		cancel()
	}()

	setLogger()

	// mongo db
	db.Initialize(ctx)

	// gin server
	engine := gin.Default()
	engine.Use(gin.Recovery())
	routes.Register(engine)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	server := http.Server{
		Addr:    ":" + port,
		Handler: engine,
	}

	if err := server.ListenAndServe(); err != nil {
		logrus.Fatalf("ListenAndServe has been failed. Error %s", err.Error())
		panic(err)
	}
}

func setLogger() {
	// logrus
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			filename := path.Base(f.File)
			return "", fmt.Sprintf("%s:%d", filename, f.Line)
		},
	})
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		gin.SetMode(gin.DebugMode)
		logrus.SetLevel(logrus.DebugLevel)
	}
	logrus.SetOutput(os.Stdout)
}

