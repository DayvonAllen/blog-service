package services

import (
	"com.aharakitchen/app/domain"
	"com.aharakitchen/app/repo"
	cache2 "github.com/go-redis/cache/v8"
)

type TagService interface {
	FindAllPostsByCategory(category, page string) (*domain.PostList, error)
	FindAllTags(rdb *cache2.Cache) (*domain.TagList, error)
}

type DefaultTagService struct {
	repo repo.TagRepo
}

func (s DefaultTagService) FindAllPostsByCategory(category, page string) (*domain.PostList, error) {
	postList, err := s.repo.FindAllPostsByCategory(category, page)
	if err != nil {
		return nil, err
	}
	return postList, nil
}

func (s DefaultTagService) FindAllTags(rdb *cache2.Cache) (*domain.TagList, error) {
	tags, err := s.repo.FindAllTags(rdb)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func NewTagService(repository repo.TagRepo) DefaultTagService {
	return DefaultTagService{repository}
}