package organization

import (
	"time"

	"github.com/estherk0/slack-ae-bot/pkg/models/base"
)

type User struct {
	base.BaseModel `bson:",inline" json:",inline"`

	UserID   string  `bson:"user_id" json:"user_id"`
	TeamID   string  `bson:"team_id" json:"team_id"`
	Name     string  `bson:"name" json:"name"`
	TZ       string  `bson:"tz" json:"tz,omitempty"`
	RealName string  `bson:"real_name" json:"real_name"`
	Medals   []Medal `bson:"medals" json:"medals"`
}

type Medal struct {
	Type      string // "gold" "silver" "bronze"
	Source    string // "karma"
	CreatedAt time.Time
}
