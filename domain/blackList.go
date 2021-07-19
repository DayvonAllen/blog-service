package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Blacklist struct {
	Id        primitive.ObjectID `bson:"_id" json:"-"`
	IP 		  string			 `bson:"ip" json:"-"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

