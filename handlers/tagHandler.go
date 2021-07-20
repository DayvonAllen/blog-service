package handlers

import (
	"com.aharakitchen/app/cache"
	"com.aharakitchen/app/domain"
	"com.aharakitchen/app/services"
	"context"
	"fmt"
	cache2 "github.com/go-redis/cache/v8"
	"github.com/gofiber/fiber/v2"
)

type TagHandler struct {
	TagService services.TagService
}

func (th *TagHandler) GetAllPostsByTags(c *fiber.Ctx) error {
	category := c.Params("category")
	page := c.Query("page", "1")

	postList, err := th.TagService.FindAllPostsByCategory(category, page)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": postList})
}

func (th *TagHandler) GetAllTags(c *fiber.Ctx) error {
	rdb := cache.RedisCachePool.Get().(*cache2.Cache)
	defer cache.RedisCachePool.Put(rdb)

	t := new(domain.TagList)
	err := rdb.Get(context.TODO(), "tagList", t)

	if err != nil {
		tags, err := th.TagService.FindAllTags(rdb)

		if err != nil {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}

		return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": tags})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": t})
}
