package karma

import (
	"context"
	"fmt"

	"github.com/estherk0/slack-ae-bot/pkg/models/karma"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack/slackevents"
	"go.mongodb.org/mongo-driver/mongo"
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
	users, err := s.karmaRepository.GetSortedUsers(ctx, currentSeason.SeasonID, 10)
	if err != nil {
		logrus.Errorln("FinishSeason failed to get top users. error: ", err.Error())
		return err
	}
	s.notifyKarmaRank(users, currentSeason.SeasonID, event.Channel)

	go func() {
		if err := s.processAwardMedals(users, fmt.Sprintf("karma_season_%d", currentSeason.SeasonID)); err != nil {
			logrus.Fatalln("Failed to award medals to rankers ", err.Error())
		}
	}()

	if err = s.karmaRepository.FinishCurrentSeason(ctx); err != nil {
		s.slackapiService.PostMessage(event.Channel, "I was trying to finish season. But I failed. Help!")
		return err
	}
	msg := fmt.Sprintf("Season #%d is finished. :partying_face:. Please check your rank and reward.", currentSeason.SeasonID)
	s.slackapiService.PostMessage(event.Channel, msg)
	return nil
}

func (s *service) processAwardMedals(rankers []karma.User, source string) error {
	ctx := context.Background()
	for i, ranker := range rankers {
		_, err := s.orgRepository.GetUser(ctx, ranker.UserID)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				slackUser, err := s.slackapiService.GetUserInfo(ranker.UserID)
				if err != nil {
					logrus.Errorf("Failed to get user info %s, index %d\n", err.Error(), i)
					return err
				}
				_, err = s.orgRepository.RegisterUser(ctx, ranker.UserID, slackUser.TeamID, slackUser.Name, slackUser.RealName)
				if err != nil {
					logrus.Errorf("Failed to register user info %s, index %d\n", err.Error(), i)
					return err
				}
			} else {
				logrus.Errorf("Failed to get user info %s, index %d\n", err.Error(), i)
				return err
			}
		}
		if err = s.orgRepository.AwardMedalToUser(ctx, ranker.UserID, source, getMedalType(i)); err != nil {
			logrus.Errorf("Failed to award medal to user %s, index %d\n", err.Error(), i)
			return err
		}

		if i == 2 { // end of process. i == 2 is bronze medal
			return nil
		}
	}
	return nil
}

func getMedalType(rank int) string {
	switch rank {
	case 0:
		return "gold"
	case 1:
		return "silver"
	case 2:
		return "bronze"
	default:
		return "stone"
	}
}
