package command

import (
	"github.com/estherk0/slack-ae-bot/pkg/controllers/command"
	"github.com/gin-gonic/gin"
)

func Register(parentGroup *gin.RouterGroup, controller command.Controller) {
	if controller == nil {
		controller = command.CreateController()
	}
	eventsGroup := parentGroup.Group("/slash")
	eventsGroup.Use()
	{
		eventsGroup.POST("", controller.HandleCommands)
	}
}
