package karma

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/estherk0/slack-ae-bot/pkg/repositories/karma"
	"github.com/estherk0/slack-ae-bot/pkg/services/slackapi"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack/slackevents"
)

const (
	receiverKarma = 1.0
	giverKarma    = 0.3
)

type Service interface {
	AddUserKarma(event *slackevents.MessageEvent) error
}

type service struct {
	slackapiService slackapi.Service
	karmaRepository karma.Repository
}

func NewService(slackapiService slackapi.Service, karmaRepository karma.Repository) Service {
	return &service{
		slackapiService: slackapiService,
		karmaRepository: karmaRepository,
	}
}

func CreateService() Service {
	return NewService(
		slackapi.CreateService(),
		karma.CreateRepository(),
	)
}

// Add only 1 karma to receiver
func (s *service) AddUserKarma(event *slackevents.MessageEvent) error {
	totalReceiverCount := 0
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	season, err := s.karmaRepository.GetCurrentSeason(ctx)
	if err != nil {
		logrus.Errorln("AddUserKarma error: ", err.Error())
	}

	giverID := event.User
	resultMessage := ""
	text := strings.ReplaceAll(event.Text, " ", "")
	r, _ := regexp.Compile(`<@\w+>\+\+`)
	matches := r.FindAllString(text, -1)
	for _, m := range matches {
		receiverID := m[2 : len(m)-3]
		if receiverID == giverID {
			resultMessage += fmt.Sprintf("Sorry, <@%s>. You are not allowed to give karma yourself.\n", giverID)
			continue
		}
		if err := s.karmaRepository.AddUserKarma(ctx, season.SeasonID, receiverID, receiverKarma); err != nil {
			return err
		}
		resultMessage += fmt.Sprintf("<@%s> has gained %0.1f karam.\n", receiverID, receiverKarma)
		totalReceiverCount += 1
	}
	if totalReceiverCount != 0 {
		if err = s.karmaRepository.AddUserKarma(ctx, season.SeasonID, giverID, giverKarma); err != nil {
			logrus.Error("failed to add point to giver ", giverID)
		} else {
			resultMessage += fmt.Sprintf("<@%s> has gained %0.1f karma.\n", giverID, giverKarma)
		}
	}
	s.slackapiService.PostMessage(event.Channel, resultMessage)
	return nil
}
