package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Tag struct {
	Id             primitive.ObjectID `bson:"_id" json:"id"`
	Value          string `bson:"value" json:"value"`
	AssociatedPosts []primitive.ObjectID `bson:"associatedPosts"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type TagDto struct {
	Id             primitive.ObjectID `bson:"_id" json:"-"`
	Value          string `bson:"value" json:"value"`
	AssociatedPosts []primitive.ObjectID `bson:"associatedPosts" json:"-"`
}

type TagList struct {
	Tags				[]TagDto `bson:"tags" json:"tags"`
	NumberOfCategories int `bson:"numberOfCategories" json:"numberOfCategories"`
}

type RedisTagList struct {
	Tags				[]byte `bson:"tags" json:"tags"`
	NumberOfCategories int `bson:"numberOfCategories" json:"numberOfCategories"`
}
