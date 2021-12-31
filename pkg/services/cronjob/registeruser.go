package cronjob

import (
	"context"

	"github.com/sirupsen/logrus"
)

func (s *service) registerUser() {
	ctx := context.Background()
	season, err := s.karmaRepository.GetCurrentSeason(ctx)
	if err != nil || season == nil {
		logrus.Warnln("[CRON][RegisterUser] failed to get current season ", err.Error())
		return
	}
	karmausers, err := s.karmaRepository.GetUsers(ctx, season.SeasonID)
	if err != nil {
		logrus.Errorln("[CRON][RegisterUser] failed to get user list ", err.Error())
		return
	}
	for _, karmaUser := range karmausers {
		orgUser, err := s.orgRepository.GetUser(ctx, karmaUser.UserID)
		if err != nil {
			logrus.Errorf("[CRON][RegisterUser] get org_users error %s for User %s\n", err.Error(), karmaUser.UserID)
			continue
		}
		if orgUser == nil { // fetch slack user info and create new user if user is not created yet,
			slackUser, err := s.slackapiService.GetUserInfo(karmaUser.UserID)
			if err != nil {
				logrus.Errorf("[CRON][RegisterUser] slackapi getuserinfo failed due to %s for User %s\n", err.Error(), karmaUser.UserID)
				continue
			}
			if _, err = s.orgRepository.RegisterUser(ctx, slackUser.ID, slackUser.TeamID, slackUser.Name, slackUser.RealName); err != nil {
				logrus.Errorf("[CRON][RegisterUser] slackapi getuserinfo failed due to %s for User %s\n", err.Error(), karmaUser.UserID)
				continue
			}
		}
	}
}
