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
	botID         = "B02QA5TMQDA"
)

type Service interface {
	AddUserKarma(event *slackevents.MessageEvent) error
	StartSeason(event *slackevents.AppMentionEvent) error
	FinishSeason(event *slackevents.AppMentionEvent) error
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
		} else if receiverID == botID {
			resultMessage += fmt.Sprintf("I don't need karma, <@%s>. But I appreciate the thought.\n", giverID)
			continue
		}
		logrus.Debugf("Receiver ID !!!", receiverID)
		if err := s.karmaRepository.AddUserKarma(ctx, season.SeasonID, receiverID, receiverKarma); err != nil {
			return err
		}
		resultMessage += fmt.Sprintf("<@%s> has gained %0.1f karma.\n", receiverID, receiverKarma)
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

func (s *service) StartSeason(event *slackevents.AppMentionEvent) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	seasonID, err := s.karmaRepository.StartNewSeason(ctx)
	if err != nil {
		msg := fmt.Sprintf("Failed to start new season because of error: %s", err.Error())
		s.slackapiService.PostMessage(event.Channel, msg)
		return err
	}
	msg := fmt.Sprintf(":tada: New Season #%d started! :tada:", seasonID)
	s.slackapiService.PostMessage(event.Channel, msg)
	return nil
}

func (s *service) FinishSeason(event *slackevents.AppMentionEvent) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	currentSeason, err := s.karmaRepository.GetCurrentSeason(ctx)
	if err != nil {
		s.slackapiService.PostMessage(event.Channel, "I couldn't find active season. :persevere: Am I wrong?")
		return err
	}
	if err = s.karmaRepository.FinishCurrentSeason(ctx); err != nil {
		s.slackapiService.PostMessage(event.Channel, "I was trying to finish season. But I failed. Help!")
		return err
	}
	msg := fmt.Sprintf("Season #%d is finished. :partying_face:. Please check your rank and reward.", currentSeason.SeasonID)
	s.slackapiService.PostMessage(event.Channel, msg)
	return nil
}
