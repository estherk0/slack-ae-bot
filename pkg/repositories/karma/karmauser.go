package karma

import (
	"context"
	"log"

	"github.com/estherk0/slack-ae-bot/pkg/models/karma"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		} else {
			return res.Err()
		}
	}

	return nil
}
func (r *repository) GetKarmaOfUser(ctx context.Context, seasonID int, userID string) (float64, error) {
	res := r.userCollection.FindOne(ctx,
		bson.M{
			"season_id": seasonID,
			"user_id":   userID,
		})
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return 0, nil
		} else {
			return -1, res.Err()
		}
	}
	user := new(karma.User)
	if err := res.Decode(&user); err != nil {
		return -1, err
	}
	return user.Karma, nil
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

func (r *repository) GetSortedUsers(ctx context.Context, seasonID int, limit int64) ([]karma.User, error) {
	findOptions := options.FindOptions{
		Limit: &limit,
		Sort: bson.M{
			"karma": -1,
		},
	}
	cursor, err := r.userCollection.Find(ctx,
		bson.M{
			"season_id": seasonID,
		},
		&findOptions,
	)
	if err != nil {
		logrus.Fatalln("GetSortedUser fatal error: ", err.Error())
		return nil, err
	}

	var users []karma.User
	if err = cursor.All(ctx, &users); err != nil {
		log.Fatalln("GetSortedUser failed to decode user: ", err.Error())
		return nil, err
	}
	return users, nil
}

func (r *repository) GetUsers(ctx context.Context, seasonID int) ([]karma.User, error) {
	cursor, err := r.userCollection.Find(ctx,
		bson.M{
			"season_id": seasonID,
		})
	if err != nil {
		return nil, err
	}
	var users []karma.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}
