package command

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/estherk0/slack-ae-bot/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

//go:generate mockery --name=Controller
type Controller interface {
	HandleCommands(c *gin.Context)
}

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

func (ctrl *controller) HandleCommands(c *gin.Context) {
	logrus.Debugln("command request arrived!")
	verifier, err := slack.NewSecretsVerifier(c.Request.Header, config.SlackSigningSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.Request.Body = ioutil.NopCloser(io.TeeReader(c.Request.Body, &verifier))
	s, err := slack.SlashCommandParse(c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	if err = verifier.Ensure(); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	switch s.Command {
	case "/echo":
		params := &slack.Msg{Text: s.Text, ResponseType: "in_channel"}
		//b, err := json.Marshal(params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		c.Writer.Header().Set("Content-Type", "application/json")
		c.JSON(http.StatusOK, params)
	default:
		logrus.Warn("Unsupported command! ", s.Command)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
}
