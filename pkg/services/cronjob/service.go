package cronjob

import (
	"github.com/estherk0/slack-ae-bot/pkg/repositories/karma"
	"github.com/estherk0/slack-ae-bot/pkg/repositories/organization"
	"github.com/estherk0/slack-ae-bot/pkg/services/slackapi"
	"github.com/robfig/cron"
)

type Service interface {
	Start()
	Stop()
}

type service struct {
	cron            *cron.Cron
	karmaRepository karma.Repository
	orgRepository   organization.Repository
	slackapiService slackapi.Service
}

func NewService(cron *cron.Cron,
	karmaRepository karma.Repository,
	orgRepository organization.Repository,
	slackapiService slackapi.Service) Service {
	return &service{
		cron: cron,
	}
}

func CreateService() Service {
	return NewService(
		cron.New(),
		karma.CreateRepository(),
		organization.CreateRepository(),
		slackapi.CreateService(),
	)
}

func (s *service) Start() {
	s.cron.AddFunc("@every 2h", s.registerUser)
	s.cron.Start()
}

func (s *service) Stop() {
	s.cron.Stop()
}
