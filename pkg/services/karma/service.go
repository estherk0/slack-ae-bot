package karma

import (
	"github.com/estherk0/slack-ae-bot/pkg/repositories/karma"
	"github.com/estherk0/slack-ae-bot/pkg/repositories/organization"
	"github.com/estherk0/slack-ae-bot/pkg/services/slackapi"
	"github.com/slack-go/slack/slackevents"
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
	orgRepository   organization.Repository
}

func NewService(
	slackapiService slackapi.Service,
	karmaRepository karma.Repository,
	orgRepository organization.Repository) Service {
	return &service{
		slackapiService: slackapiService,
		karmaRepository: karmaRepository,
		orgRepository:   orgRepository,
	}
}

func CreateService() Service {
	return NewService(
		slackapi.CreateService(),
		karma.CreateRepository(),
		organization.CreateRepository(),
	)
}
