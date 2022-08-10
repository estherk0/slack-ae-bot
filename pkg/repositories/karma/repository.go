package karma

import (
	"context"

	"github.com/estherk0/slack-ae-bot/pkg/db"
	"github.com/estherk0/slack-ae-bot/pkg/models/karma"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	AddUserKarma(ctx context.Context, seasonID int, userID string, karma float64) error
	GetKarmaOfUser(ctx context.Context, seasonID int, userID string) (float64, error)
	GetUsers(ctx context.Context, seasonID int) ([]karma.User, error)
	GetCurrentSeason(ctx context.Context) (*karma.Season, error)
	StartNewSeason(ctx context.Context) (int64, error)
	FinishCurrentSeason(ctx context.Context) error
	GetSortedUsers(ctx context.Context, seasonID int, limit int64) ([]karma.User, error)
	CreateNewLog(ctx context.Context, seasonID int, receiverID string, giverID string, eventTimestamp string)
	SearchLogs(ctx context.Context, seasonID int, receiverID string, days int) ([]karma.Log, error)
	GetLog(ctx context.Context, seasonID int, receiverID string, giverID string, eventTimestamp string) (*karma.Log, error)
}

type repository struct {
	userCollection   *mongo.Collection
	seasonCollection *mongo.Collection
	logCollection    *mongo.Collection
}

func NewRepository(userCollection *mongo.Collection,
	seasonCollection *mongo.Collection,
	logCollection *mongo.Collection) Repository {
	return &repository{
		userCollection:   userCollection,
		seasonCollection: seasonCollection,
		logCollection:    logCollection,
	}
}

func CreateRepository() Repository {
	return NewRepository(
		db.GetDB().Collection(userCollectionName),
		db.GetDB().Collection(seasonCollectionName),
		db.GetDB().Collection(logCollectionName),
	)
}
