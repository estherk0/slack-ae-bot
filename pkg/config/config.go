package config

import (
	"os"
)

var (
	SlackSigningSecret string
	SlackToken         string
	DBHost             string
	DBPort             string
	DBUserName         string
	DBPassword         string
	DBName             string
)

func init() {
	SlackSigningSecret = os.Getenv("SLACK_SIGNING_SECRET")
	SlackToken = os.Getenv("SLACK_TOKEN")
	DBHost = os.Getenv("DB_HOST")
	if DBHost == "" {
		DBHost = "localhost"
	}
	DBPort = os.Getenv("DB_PORT")
	if DBPort == "" {
		DBPort = "27017"
	}
	DBUserName = os.Getenv("DB_USERNAME")
	DBPassword = os.Getenv("DB_PASSWORD")
	DBName = os.Getenv("DB_NAME")
}
