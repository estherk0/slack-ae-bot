package karma

import (
	"time"

	"github.com/estherk0/slack-ae-bot/pkg/models/base"
)

type Log struct {
	base.BaseModel `bson:",inline" json:",inline"`

	SeasonID       int       `bson:"season_id" json:"season_id"`
	ReceiverID     string    `bson:"receiver_id" json:"receiver_id"`
	GiverID        string    `bson:"giver_id" json:"giver_id"`
	CreatedAt      time.Time `bson:"created_at" json:"created_at"`
	EventTimestamp string    `bson:"event_timestamp" json:"event_timestamp"`
}
