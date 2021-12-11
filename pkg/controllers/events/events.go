package events

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/estherk0/slack-ae-bot/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

//go:generate mockery --name=Controller
type Controller interface {
	HandleEvents(c *gin.Context)
}

var (
	api = slack.New(config.SlackToken)
)

type controller struct {
}

// NewController -
func NewController() *controller {
	return &controller{}
}

// CreateController -
func CreateController() *controller {
	return NewController()
}

// Handler -
func (ctrl *controller) HandleEvents(c *gin.Context) {
	logrus.Debugln("event arrived!")
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	sv, err := slack.NewSecretsVerifier(c.Request.Header, config.SlackSigningSecret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	if _, err := sv.Write(body); err != nil {
		logrus.Debugln("events body write failed!")
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	if err := sv.Ensure(); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}
	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			api.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
		case *slackevents.MessageEvent:
			api.PostMessage(ev.Channel, slack.MsgOptionText("I am listening your message.", false))
		}
	}
	c.JSON(http.StatusOK, gin.H{})
}
