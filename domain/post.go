package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Post struct {
	Id        primitive.ObjectID `bson:"_id" json:"-"`
	Title     string             `bson:"title" json:"title"`
	Content   string             `bson:"content" json:"content"`
	Preview   string             `bson:"preview" json:"preview"`
	Author    string             `bson:"author" json:"author"`
	MainImage string             `bson:"mainImage" json:"mainImage"`
	StoryImages []string		 `bson:"storyImages" json:"storyImages"`
	Score     int                `bson:"score" json:"-"`
	Tag      string           `bson:"tag" json:"tag"`
	Visible	  bool				 `bson:"visible" json:"visible"`
	Updated   bool               `bson:"updated" json:"updated"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type PostDto struct {
	Title     string             `bson:"title" json:"title"`
	Content   string             `bson:"content" json:"content"`
	Author    string             `bson:"author" json:"author"`
	MainImage string             `bson:"mainImage" json:"mainImage"`
	StoryImages []string		 `bson:"storyImages" json:"storyImages"`
	Tag      string           `bson:"tag" json:"tag"`
	Updated   bool               `bson:"updated" json:"updated"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type RedisPostDto struct {
	Title     string             `bson:"title" json:"title"`
	Content   string             `bson:"content" json:"content"`
	Author    string             `bson:"author" json:"author"`
	MainImage string             `bson:"mainImage" json:"mainImage"`
	StoryImages []byte		 `bson:"storyImages" json:"storyImages"`
	Tag      string           `bson:"tag" json:"tag"`
	Updated   bool               `bson:"updated" json:"updated"`
	CreatedAt []byte         `bson:"createdAt" json:"createdAt"`
	UpdatedAt []byte          `bson:"updatedAt" json:"updatedAt"`
}

type PostPreviewDto struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	Title     string             `json:"title"`
	Preview   string             `json:"preview"`
	Author    string             `json:"author"`
	MainImage string             `json:"mainImage"`
	Tag      string           `json:"tag"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
}

type PostList struct {
	Posts []PostPreviewDto `json:"posts"`
	NumberOfPosts int64		`json:"numberOfPosts"`
	CurrentPage int			`json:"currentPage"`
	NumberOfPages int		`json:"numberOfPages"`
}

type RedisPostList struct {
	Posts []byte `json:"posts"`
	NumberOfPosts int64		`json:"numberOfPosts"`
	CurrentPage int			`json:"currentPage"`
	NumberOfPages int		`json:"numberOfPages"`
}