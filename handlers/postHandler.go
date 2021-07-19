package handlers

import (
	"com.aharakitchen/app/cache"
	"com.aharakitchen/app/domain"
	"com.aharakitchen/app/services"
	"context"
	"fmt"
	cache2 "github.com/go-redis/cache/v8"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
)

type PostHandler struct {
	PostService services.PostService
}

func (ph *PostHandler) GetAllPosts(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	newStoriesQuery := c.Query("new", "false")

	isNew, err := strconv.ParseBool(newStoriesQuery)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("must provide a valid value")})
	}

	postList, err := ph.PostService.FindAllPosts(page, isNew)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": postList})
}

func (ph *PostHandler) GetFeaturedPosts(c *fiber.Ctx) error {
	rdb := cache.RedisCachePool.Get().(*cache2.Cache)
	defer cache.RedisCachePool.Put(rdb)

	pl := new(domain.PostList)
	err := rdb.Get(context.TODO(), "featuredstories", pl)

	if err != nil {
		postList, err := ph.PostService.FeaturedPosts(rdb)

		if err != nil {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}

		return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": postList})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": pl})
}

func (ph *PostHandler) GetPostById(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	post, err := ph.PostService.FindPostById(id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": post})
}

