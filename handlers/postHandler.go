package handlers

import (
	"com.aharakitchen/app/database"
	"com.aharakitchen/app/domain"
	"com.aharakitchen/app/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/vmihailenco/msgpack/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"time"
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
	rdb := database.Conn.Get()

	v, err := redis.Values(rdb.Do("HGETALL", "featuredstories"))

	if err != nil || len(v) == 0 {
		postList, err := ph.PostService.FeaturedPosts()

		if err != nil {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}

		return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": postList})
	}

	rpl := new(domain.RedisPostList)

	if err := redis.ScanStruct(v, rpl); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	var posts []domain.PostPreviewDto

	err = msgpack.Unmarshal(rpl.Posts, &posts)

	pl := new(domain.PostList)

	pl.Posts = posts
	pl.NumberOfPosts = rpl.NumberOfPosts
	pl.NumberOfPages = rpl.NumberOfPages
	pl.CurrentPage = rpl.CurrentPage

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": pl})
}

func (ph *PostHandler) GetPostById(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	rdb := database.Conn.Get()

	r := new(domain.RedisPostDto)

	v, err := redis.Values(rdb.Do("HGETALL", id.String()+"getbyID"))

	if len(v) == 0 || err != nil {
		post, err := ph.PostService.FindPostById(id)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
		}
		return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": post})
	}

	if err := redis.ScanStruct(v, r); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	p := new(domain.PostDto)

	var urls []string
	var created time.Time
	var updated time.Time

	err = msgpack.Unmarshal(r.StoryImages, &urls)


	if err != nil {
		fmt.Println(err)
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = msgpack.Unmarshal(r.CreatedAt, &created)

	if err != nil {
		fmt.Println(err)
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	err = msgpack.Unmarshal(r.UpdatedAt, &updated)

	if err != nil {
		fmt.Println(err)
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "error...", "data": fmt.Sprintf("%v", err)})
	}

	p.Title= r.Title
	p.Tag = r.Tag
	p.Content = r.Content
	p.MainImage= r.MainImage
	p.Author= r.Author
	p.CreatedAt= created
	p.UpdatedAt= updated
	p.Updated= r.Updated
	p.StoryImages= urls

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "success", "data": p})
}
