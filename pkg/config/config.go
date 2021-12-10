package config

import (
	"os"
)

var (
	SlackSigningSecret string
	SlackToken         string
)

func init() {
	SlackSigningSecret = os.Getenv("SLACK_SIGNING_SECRET")
	SlackToken = os.Getenv("SLACK_TOKEN")
}
