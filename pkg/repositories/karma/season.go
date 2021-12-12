package karma

import (
	"context"
	"errors"
	"time"

	"github.com/estherk0/slack-ae-bot/pkg/models/karma"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

func (r *repository) GetCurrentSeason(ctx context.Context) (*karma.Season, error) {
	res := r.seasonCollection.FindOne(ctx, bson.M{"in_progress": true})
	if res.Err() != nil {
		return nil, res.Err()
	}
	season := new(karma.Season)
	if err := res.Decode(&season); err != nil {
		return nil, err
	}
	return season, nil
}

func (r *repository) getSeasonTotalCount(ctx context.Context) (int64, error) {
	return r.seasonCollection.CountDocuments(ctx, bson.M{})
}
func (r *repository) FinishCurrentSeason(ctx context.Context) error {
	res, err := r.seasonCollection.UpdateOne(ctx,
		bson.M{"in_progress": true},
		bson.M{
			"in_progress": false,
			"finished_at": time.Now(),
		},
	)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.New("no season finished")
	}
	return nil
}

func (r *repository) StartNewSeason(ctx context.Context) (int64, error) {
	season, _ := r.GetCurrentSeason(ctx)
	if season != nil {
		return -1, errors.New("please finish current season first")
	}
	count, err := r.getSeasonTotalCount(ctx)
	if err != nil {
		logrus.Error("failed to get season total count ", err.Error())
		return -1, err
	}
	newSeasonID := count + 1
	_, err = r.seasonCollection.InsertOne(ctx,
		bson.M{
			"season_id":   newSeasonID,
			"in_progress": true,
		})
	if err != nil {
		return -1, err
	}
	logrus.Info("new season started season id: ", season.SeasonID)
	return newSeasonID, nil
}
