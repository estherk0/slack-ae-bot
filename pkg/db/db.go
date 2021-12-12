package db

import (
	"context"
	"fmt"

	"github.com/estherk0/slack-ae-bot/pkg/config"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client   *mongo.Client
	database *mongo.Database
)

func GetDB() *mongo.Database {
	return database
}

func Initialize(ctx context.Context) (err error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", config.DBUserName, config.DBPassword, config.DBHost, config.DBPort)
	logrus.Debugln("Trying to connect %s ...", uri)
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	database = client.Database(config.DBName)
	return nil
}

func Close(ctx context.Context) error {
	return client.Disconnect(ctx)
}
