package config

import (
	"os"
)

var (
	SlackSigningSecret = os.Getenv("SLACK_SIGNING_SECRET")
	SlackToken         = os.Getenv("SLACK_TOKEN")
)
