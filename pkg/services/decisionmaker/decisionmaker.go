package decisionmaker

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/estherk0/slack-ae-bot/pkg/services/slackapi"
	"github.com/slack-go/slack/slackevents"
)

var resultFmtStr = [...]string{
	"Do it now! <@%s> :fire:",
	"Never ever ever <@%s> :smiling_imp:",
	"Hmm.... Well... <@%s> :thinking_face:",
	"Yes please. <@%s> :blobnodfast:",
}

type Service interface {
	MakeDecision(event *slackevents.AppMentionEvent)
}

type service struct {
	slackapiService slackapi.Service
}

func NewService(slackapiService slackapi.Service) Service {
	return &service{
		slackapiService: slackapiService,
	}
}

func CreateService() Service {
	return NewService(slackapi.CreateService())
}

func (s *service) MakeDecision(event *slackevents.AppMentionEvent) {
	decisionCount := len(resultFmtStr)
	randSource := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(randSource)
	idx := r1.Intn(100) % decisionCount
	msg := fmt.Sprintf(resultFmtStr[idx], event.User)
	s.slackapiService.PostMessage(event.Channel, msg)
}
