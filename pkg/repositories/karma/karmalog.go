package karma

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/estherk0/slack-ae-bot/pkg/models/karma"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

const (
	logCollectionName = "karma_log"
)

func (r *repository) CreateNewLog(ctx context.Context, seasonID int, receiverID string, giverID string, eventTimestamp string) {
	res, err := r.logCollection.InsertOne(ctx,
		bson.M{
			"season_id":       seasonID,
			"receiver_id":     receiverID,
			"giver_id":        giverID,
			"created_at":      time.Now(),
			"event_timestamp": eventTimestamp,
		},
	)
	if err != nil {
		logrus.Errorf("[CreateNewLog] failed to create karma log. error: %s", err.Error())
	}
	logrus.Info("[CreateNewLog] new karma_log created. id: ", res.InsertedID)
}

func (r *repository) GetLog(ctx context.Context, seasonID int, receiverID string, giverID string, eventTimestamp string) (*karma.Log, error) {

	res := r.logCollection.FindOne(ctx, bson.M{
		"season_id":       seasonID,
		"receiver_id":     receiverID,
		"giver_id":        giverID,
		"event_timestamp": eventTimestamp,
	})

	if res == nil || res.Err() != nil {
		return nil, fmt.Errorf("[SearchLogs] failed to find karma log. error: %s", res.Err().Error())
	}

	log := new(karma.Log)
	if err := res.Decode(&log); err != nil {
		return nil, err
	}
	return log, nil
}

func (r *repository) SearchLogs(ctx context.Context, seasonID int, receiverID string, days int) ([]karma.Log, error) {
	findOptions := options.FindOptions{
		Sort: bson.M{
			"created_at": -1,
		},
	}
	date := time.Now().AddDate(0, 0, -days)
	cursor, err := r.logCollection.Find(ctx, bson.M{
		"season_id":   seasonID,
		"receiver_id": receiverID,
		"created_at": bson.M{
			"$gte": date,
		},
	}, &findOptions)

	if err != nil {
		return nil, fmt.Errorf("[SearchLogs] failed to find karma logs. error: %s", err.Error())
	}

	var karmaLogs []karma.Log
	if err = cursor.All(ctx, &karmaLogs); err != nil {
		log.Fatalln("[SearchLogs] failed to decode karma logs: ", err.Error())
		return nil, err
	}
	return karmaLogs, nil
}
