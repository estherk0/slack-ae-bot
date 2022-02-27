package events

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/estherk0/slack-ae-bot/pkg/config"
	"github.com/estherk0/slack-ae-bot/pkg/services/decisionmaker"
	"github.com/estherk0/slack-ae-bot/pkg/services/karma"
	"github.com/estherk0/slack-ae-bot/pkg/services/slackapi"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

const karmaMessagePattern = `<@\w+>\s*\+\+`

var randomResponse = [...]string{
	"Sorry, I don't understand what you are saying. :sob:",
	"Don't bother me. :blobdizzy:",
	"Please don't say anything more. :no_mouth:",
	"I am trying to understand your language.",
	"Naevis, calling.",
	"Don't you know I am a savage? :sunglasses:",
	"(zu zu zu zu)",
}

//go:generate mockery --name=Controller
type Controller interface {
	HandleEvents(c *gin.Context)
}

type controller struct {
	karmaService         karma.Service
	slackapiService      slackapi.Service
	decisionmakerService decisionmaker.Service
}

// NewController -
func NewController(karmaService karma.Service,
	slackapiService slackapi.Service,
	decisionmakerService decisionmaker.Service) *controller {
	return &controller{
		karmaService:         karmaService,
		slackapiService:      slackapiService,
		decisionmakerService: decisionmakerService,
	}
}

// CreateController -
func CreateController() *controller {
	return NewController(
		karma.CreateService(),
		slackapi.CreateService(),
		decisionmaker.CreateService(),
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
	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		c.Writer.Header().Set("Content-Type", "text")
		c.JSON(http.StatusOK, r.Challenge)
		return
	} else if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			ctrl.appMentionEvent(ev)
		case *slackevents.MessageEvent:
			ctrl.messageEvent(ev)
		}
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	logrus.Errorf("unknown event api type %s", eventsAPIEvent.Type)
	c.JSON(http.StatusInternalServerError, gin.H{})
}

func (ctrl *controller) appMentionEvent(event *slackevents.AppMentionEvent) {
	if strings.Contains(event.Text, "karma") { // karma
		if strings.Contains(event.Text, "season") {
			if strings.Contains(event.Text, "start") {
				ctrl.karmaService.StartSeason(event)
			} else if strings.Contains(event.Text, "finish") || strings.Contains(event.Text, "end") {
				ctrl.karmaService.FinishSeason(event)
			}
		} else if strings.Contains(event.Text, "my") { // query current karma point
			if err := ctrl.karmaService.GetKarmaOfUser(event); err != nil {
				logrus.Errorf("GetKarmaOfUser error %s", err.Error())
			}
		} else if strings.Contains(event.Text, "history") { // query karma history
			if err := ctrl.karmaService.GetHistories(event); err != nil {
				logrus.Errorf("GetHistories error %s", err.Error())
			}
		} else if strings.Contains(event.Text, "top") { // karma current top list
			ctrl.karmaService.GetTopKarmaUsers(event)
		}
	} else if strings.Contains(event.Text, "??") { // decision maker
		ctrl.decisionmakerService.MakeDecision(event)
	} else {
		ctrl.unknownCommandResponse(event)
	}
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

func (ctrl *controller) unknownCommandResponse(event *slackevents.AppMentionEvent) {
	decisionCount := len(randomResponse)
	randSource := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(randSource)
	idx := r1.Intn(10000) % decisionCount
	ctrl.slackapiService.PostMessage(event.Channel, randomResponse[idx])
}
