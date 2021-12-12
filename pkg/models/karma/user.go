package karma

import "github.com/estherk0/slack-ae-bot/pkg/models/base"

type User struct {
	base.BaseModel `bson:",inline" json:",inline"`

	UserID   string  `bson:"user_id" json:"user_id"`
	SeasonID int     `bson:"season_id" json:"season_id"`
	Karma    float64 `bson:"karma" json:"karma"`
}
