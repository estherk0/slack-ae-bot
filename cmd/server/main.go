package main

import (
	"net/http"
	"os"

	"github.com/estherk0/slack-ae-bot/pkg/routes"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
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

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logrus.Fatalf("ListenAndServe has been failed. Error %s", err.Error())
			panic(err)
		}
	}()
}
