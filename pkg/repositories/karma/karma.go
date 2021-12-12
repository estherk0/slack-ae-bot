package karma

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

const (
	userCollectionName   = "karma_user"
	seasonCollectionName = "karma_season"
)

func (r *repository) AddUserKarma(ctx context.Context, seasonID int, userID string, karma float64) error {
	res := r.userCollection.FindOneAndUpdate(ctx,
		bson.M{
			"season_id": seasonID,
			"user_id":   userID,
		},
		bson.M{"$inc": bson.M{"karma": karma}})
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			if err := r.createUser(ctx, seasonID, userID, karma); err != nil {
				logrus.Errorf("failed to create user %s, error: ", userID, err.Error())
				return err
			}
		}
		return res.Err()
	}

	return nil
}

func (r *repository) createUser(ctx context.Context, seasonID int, userID string, karma float64) error {
	res, err := r.userCollection.InsertOne(ctx,
		bson.M{
			"season_id": seasonID,
			"user_id":   userID,
			"karma":     karma,
		})
	if err != nil {
		return err
	}
	logrus.Info("new karma_user created. id: ", res.InsertedID)
	return nil
}
