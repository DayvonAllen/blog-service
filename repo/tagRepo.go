package repo

import (
	"com.aharakitchen/app/domain"
	cache2 "github.com/go-redis/cache/v8"
)

type TagRepo interface {
	FindAllPostsByCategory(category, page string) (*domain.PostList, error)
	Create(tag domain.Tag) error
	FindAllTags(rdb *cache2.Cache) (*domain.TagList, error)
	UpdateTag(tag domain.Tag) error
}
