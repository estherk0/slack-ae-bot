package events

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/estherk0/slack-ae-bot/pkg/config"
	"github.com/estherk0/slack-ae-bot/pkg/services/karma"
	"github.com/estherk0/slack-ae-bot/pkg/services/slackapi"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

const karmaMessagePattern = `<@\w+>\s*\+\+`

//go:generate mockery --name=Controller
type Controller interface {
	HandleEvents(c *gin.Context)
}

type controller struct {
	karmaService    karma.Service
	slackapiService slackapi.Service
}

// NewController -
func NewController(karmaService karma.Service, slackapiService slackapi.Service) *controller {
	return &controller{
		karmaService:    karmaService,
		slackapiService: slackapiService,
	}
}

// CreateController -
func CreateController() *controller {
	return NewController(
		karma.CreateService(),
		slackapi.CreateService(),
	)
}

// Handler -
func (ctrl *controller) HandleEvents(c *gin.Context) {
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
			ctrl.appMentionEvent(ev)
		case *slackevents.MessageEvent:
			ctrl.messageEvent(ev)
		}
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (ctrl *controller) appMentionEvent(event *slackevents.AppMentionEvent) {
	ctrl.slackapiService.PostMessage(event.Channel, "Yes, hello.")
}

func (ctrl *controller) messageEvent(event *slackevents.MessageEvent) {
	logrus.Debugln("Message Event received! ", event.BotID, event.User, event.Channel)
	if event.BotID != "" { // Ignore bot message
		return
	}
	re, _ := regexp.Compile(karmaMessagePattern)
	if re.MatchString(event.Text) {
		err := ctrl.karmaService.AddUserKarma(event)
		if err != nil {
			logrus.Error("add user karma has failed due to error: ", err.Error())
		}
	}
}
