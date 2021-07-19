package repo

import (
	"com.aharakitchen/app/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostRepo interface {
	FindAllPosts(page string, newPosts bool) (*domain.PostList, error)
	FeaturedPosts() (*domain.PostList, error)
	Create(post domain.Post) error
	UpdateByTitle(post domain.Post) error
	FindPostById(id primitive.ObjectID) (*domain.PostDto, error)
}
