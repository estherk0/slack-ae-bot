package karma

import (
	"context"
	"errors"

	"github.com/estherk0/slack-ae-bot/pkg/models/karma"
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

func (r *repository) FinishCurrentSeason(ctx context.Context) error {
	res, err := r.seasonCollection.UpdateOne(ctx, bson.M{"in_progress": true}, bson.M{"in_progress": false})
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.New("no season finished")
	}
	return nil
}
