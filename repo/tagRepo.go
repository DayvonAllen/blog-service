package repo

import (
	"com.aharakitchen/app/domain"
)

type TagRepo interface {
	FindAllPostsByCategory(category, page string) (*domain.PostList, error)
	Create(tag domain.Tag) error
	FindAllTags() (*domain.TagList, error)
	UpdateTag(tag domain.Tag) error
}
