package karma

import "github.com/slack-go/slack/slackevents"

const (
	receiverPoint = 1
	giverPoint    = 0.3
)

type Service interface {
	AddUserKarma(event *slackevents.MessageEvent) error
}

type service struct {
}

func NewService() Service {
	return &service{}
}

func CreateService() Service {
	return NewService()
}

// Add only 1 karma to receiver
func (svc *service) AddUserKarma(event *slackevents.MessageEvent) error {
	return nil
}
