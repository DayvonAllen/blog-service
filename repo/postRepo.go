package repo

import (
	"com.aharakitchen/app/domain"
	cache2 "github.com/go-redis/cache/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostRepo interface {
	FindAllPosts(page string, newPosts bool) (*domain.PostList, error)
	FeaturedPosts(cache *cache2.Cache) (*domain.PostList, error)
	Create(post domain.Post) error
	UpdateByTitle(post domain.Post) error
	DeleteById(post domain.Post) error
	FindPostById(id primitive.ObjectID, rdb *cache2.Cache) (*domain.PostDto, error)
}
