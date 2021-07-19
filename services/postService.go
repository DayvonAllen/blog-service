package services

import (
	"com.aharakitchen/app/domain"
	"com.aharakitchen/app/repo"
	cache2 "github.com/go-redis/cache/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostService interface {
	FindAllPosts(page string, newPosts bool) (*domain.PostList, error)
	FeaturedPosts(cache *cache2.Cache) (*domain.PostList, error)
	FindPostById(id primitive.ObjectID) (*domain.PostDto, error)
}

type DefaultPostService struct {
	repo repo.PostRepo
}

func (s DefaultPostService) FindAllPosts(page string, newPosts bool) (*domain.PostList, error) {
	postList, err := s.repo.FindAllPosts(page, newPosts)
	if err != nil {
		return nil, err
	}
	return postList, nil
}

func (s DefaultPostService) FeaturedPosts(cache *cache2.Cache) (*domain.PostList, error) {
	postList, err := s.repo.FeaturedPosts(cache)
	if err != nil {
		return nil, err
	}
	return postList, nil
}

func (s DefaultPostService) FindPostById(id primitive.ObjectID) (*domain.PostDto, error) {
	post, err := s.repo.FindPostById(id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func NewPostService(repository repo.PostRepo) DefaultPostService {
	return DefaultPostService{repository}
}