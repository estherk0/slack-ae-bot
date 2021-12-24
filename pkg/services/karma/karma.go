package karma

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	karmamodel "github.com/estherk0/slack-ae-bot/pkg/models/karma"
	"github.com/estherk0/slack-ae-bot/pkg/repositories/karma"
	"github.com/estherk0/slack-ae-bot/pkg/services/slackapi"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack/slackevents"
)

const (
	receiverKarma = 1.0
	giverKarma    = 0.3
	botID         = "U02QXJZUNC8"
)

type Service interface {
	AddUserKarma(event *slackevents.MessageEvent) error
	GetUserKarma(event *slackevents.AppMentionEvent) error
	StartSeason(event *slackevents.AppMentionEvent) error
	FinishSeason(event *slackevents.AppMentionEvent) error
	GetTopKarmaUsers(event *slackevents.AppMentionEvent) error
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

const karmaTopTmpl = `:mega: Karma Season #{{ .SeasonID }} Ranking!
{{- range $index, $user := .Users }}
    Rank {{ $index +1 }} <@{{ $user.UserID }}> Karma: {{ $user.Karma }}
{{- end }}
`

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

func (s *service) GetUserKarma(event *slackevents.AppMentionEvent) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	season, err := s.karmaRepository.GetCurrentSeason(ctx)
	if err != nil {
		logrus.Errorln("GetCurrentSeason error: ", err.Error())
	}
	karma, err := s.karmaRepository.GetUserKarma(ctx, season.SeasonID, event.User)
	if err != nil {
		return err
	}
	msg := fmt.Sprintf("<@%s>'s karma is %0.1f", event.User, karma)
	s.slackapiService.PostMessage(event.Channel, msg)
	return nil
}

func (s *service) GetTopKarmaUsers(event *slackevents.AppMentionEvent) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	season, err := s.karmaRepository.GetCurrentSeason(ctx)
	if err != nil {
		logrus.Errorln("GetTopKarmaUsers failed to get season. error: ", err.Error())
		return err
	}
	users, err := s.karmaRepository.GetSortedUsers(ctx, season.SeasonID, 10)
	if err != nil {
		logrus.Errorln("GetTopKarmaUsers failed to get top users. error: ", err.Error())
		return err
	}

	// debug
	for i, user := range users {
		logrus.Infoln("rank", i, "ID", user.ID, "karma", user.Karma)
	}
	t, err := template.New("karma top template").Parse(karmaTopTmpl)
	if err != nil {
		logrus.Errorln("GetTopKarmaUsers failed to template error ", err.Error())
		return err
	}
	var tpl bytes.Buffer
	res := struct {
		Users    []karmamodel.User
		SeasonID int
	}{
		users,
		season.SeasonID,
	}
	t.Execute(&tpl, res)
	if err = s.slackapiService.PostMessage(event.Channel, tpl.String()); err != nil {
		logrus.Errorln("GetTopKarmaUsers failed to send message. error: ", err.Error())
		return err
	}
	return nil
}
