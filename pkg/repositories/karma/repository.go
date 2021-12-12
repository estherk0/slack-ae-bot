package karma

import (
	"context"

	"github.com/estherk0/slack-ae-bot/pkg/db"
	"github.com/estherk0/slack-ae-bot/pkg/models/karma"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	AddUserKarma(ctx context.Context, seasonID int, userID string, karma float64) error
	GetCurrentSeason(ctx context.Context) (*karma.Season, error)
	StartNewSeason(ctx context.Context) (int64, error)
	FinishCurrentSeason(ctx context.Context) error
}

type repository struct {
	userCollection   *mongo.Collection
	seasonCollection *mongo.Collection
}

func NewRepository(userCollection *mongo.Collection, seasonCollection *mongo.Collection) Repository {
	return &repository{
		userCollection:   userCollection,
		seasonCollection: seasonCollection,
	}
}

func CreateRepository() Repository {
	return NewRepository(
		db.GetDB().Collection(userCollectionName),
		db.GetDB().Collection(seasonCollectionName),
	)
}
