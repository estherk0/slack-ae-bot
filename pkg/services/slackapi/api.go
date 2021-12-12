package slackapi

import (
	"github.com/estherk0/slack-ae-bot/pkg/config"
	"github.com/slack-go/slack"
)

var (
	client = slack.New(config.SlackToken)
)

type Service interface {
	PostMessage(channelID, text string) error
}

type service struct {
}

func NewService() Service {
	return &service{}
}

func CreateService() Service {
	return NewService()
}

func (svc *service) PostMessage(channelID, text string) error {
	_, _, err := client.PostMessage(channelID, slack.MsgOptionText(text, false))
	return err
}
