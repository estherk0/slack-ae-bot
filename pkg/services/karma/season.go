package karma

import (
	"context"
	"fmt"

	"github.com/slack-go/slack/slackevents"
)

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
