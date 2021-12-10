package events

import (
	"github.com/estherk0/slack-ae-bot/pkg/controllers/events"
	"github.com/gin-gonic/gin"
)

func Register(parentGroup *gin.RouterGroup, controller events.Controller) {
	if controller == nil {
		controller = events.CreateController()
	}
	eventsGroup := parentGroup.Group("/events-endpoint")
	eventsGroup.Use()
	{
		eventsGroup.POST("", controller.HandleEvents)
	}
}
