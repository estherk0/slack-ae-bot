package karma

import (
	"time"

	"github.com/estherk0/slack-ae-bot/pkg/models/base"
)

type Season struct {
	base.BaseModel `bson:",inline" json:",inline"`

	SeasonID   int       `bson:"season_id" json:"season_id"`
	FinishedAt time.Time `bson:"finished_at" json:"finished_at"`
	InProgress bool      `bson:"in_progress" json:"in_progress"`
}
