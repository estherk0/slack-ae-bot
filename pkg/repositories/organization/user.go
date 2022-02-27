package organization

import (
	"context"
	"fmt"
	"time"

	"github.com/estherk0/slack-ae-bot/pkg/models/organization"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

func (r *repository) GetUser(ctx context.Context, userID string) (*organization.User, error) {
	res := r.userCollection.FindOne(ctx, bson.M{"user_id": userID})
	if res.Err() != nil {
		return nil, res.Err()
	}
	user := new(organization.User)
	if err := res.Decode(&user); err != nil {
		logrus.Errorln("GetSlackUser failed to decode user: ", err.Error())
		return nil, err
	}
	return user, nil
}

func (r *repository) RegisterUser(ctx context.Context, userID, teamID, name, realName string) (*organization.User, error) {
	newUser := organization.User{
		UserID:   userID,
		TeamID:   teamID,
		Name:     name,
		RealName: realName,
	}
	res, err := r.userCollection.InsertOne(ctx, newUser)
	if err != nil || res.InsertedID != nil {
		return nil, err
	}
	return &newUser, nil
}

func (r *repository) AwardMedalToUser(ctx context.Context, userID, source, medalType string) error {
	medal := organization.Medal{
		Source:    source,
		Type:      medalType,
		CreatedAt: time.Now(),
	}
	matchQuery := bson.M{
		"user_id": userID,
	}
	updateQuery := bson.M{
		"$push": bson.M{
			"org_users.$.medals": medal,
		},
	}
	res, err := r.userCollection.UpdateOne(ctx, matchQuery, updateQuery)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return fmt.Errorf("[AwardMedalToUser] awardMedalToUser: matching user not found for userID %s" + userID)
	}
	return nil
}
