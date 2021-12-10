package routes

import (
	"net/http"

	"github.com/estherk0/slack-ae-bot/pkg/routes/command"
	"github.com/estherk0/slack-ae-bot/pkg/routes/events"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"response": "pong",
		})
	})
	rootGroup := r.Group("/")
	events.Register(rootGroup, nil)
	command.Register(rootGroup, nil)
}
