package organization

import (
	"context"

	"github.com/estherk0/slack-ae-bot/pkg/db"
	"github.com/estherk0/slack-ae-bot/pkg/models/organization"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	userCollectionName = "org_users"
)

type Repository interface {
	GetUser(ctx context.Context, userID string) (*organization.User, error)
	RegisterUser(ctx context.Context, userID, teamID, name, realName string) (*organization.User, error)
	AwardMedalToUser(ctx context.Context, userID, source, medalType string) error
}

type repository struct {
	userCollection *mongo.Collection
}

func NewRepository(userCollection *mongo.Collection) Repository {
	return &repository{
		userCollection: userCollection,
	}
}

func CreateRepository() Repository {
	return NewRepository(
		db.GetDB().Collection(userCollectionName),
	)
}
