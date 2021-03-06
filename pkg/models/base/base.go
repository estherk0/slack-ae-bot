package base

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type BaseModel struct {
	ID        bson.ObjectId `bson:"_id" json:"_id"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}

func (b *BaseModel) SetCreated(t time.Time) {
	b.CreatedAt = t
}

func (b *BaseModel) SetModified(t time.Time) {
	b.UpdatedAt = t
}
