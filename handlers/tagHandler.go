package handlers

import (
	"com.aharakitchen/app/database"
	"com.aharakitchen/app/domain"
	"com.aharakitchen/app/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/vmihailenco/msgpack/v5"
	"runtime"
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
	fmt.Println(runtime.NumGoroutine())

	rdb := database.Conn.Get()

	v, err := redis.Values(rdb.Do("HGETALL", "tags"))

	if err != nil || len(v) == 0 {
		tags, err := th.TagService.FindAllTags()

		if err != nil {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}

		return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": tags})
	}

	rtl := new(domain.RedisTagList)

	if err := redis.ScanStruct(v, rtl); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	t := new(domain.TagList)

	var tl []domain.TagDto

	err = msgpack.Unmarshal(rtl.Tags, &tl)

	t.Tags = tl
	t.NumberOfCategories = rtl.NumberOfCategories

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": t})
}
