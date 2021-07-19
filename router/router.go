package router

import (
	"com.aharakitchen/app/handlers"
	"com.aharakitchen/app/repo"
	"com.aharakitchen/app/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetupRoutes(app *fiber.App) {
	th := handlers.TagHandler{TagService: services.NewTagService(repo.NewTagRepoImpl())}
	ph := handlers.PostHandler{PostService: services.NewPostService(repo.NewPostRepoImpl())}

	app.Use(recover.New())
	api := app.Group("", logger.New())

	tags := api.Group("/tags")
	tags.Get("/category/:category", th.GetAllPostsByTags)
	tags.Get("/", th.GetAllTags)

	posts := api.Group("/posts")
	posts.Get("/featured", ph.GetFeaturedPosts)
	posts.Get("/:id", ph.GetPostById)
	posts.Get("/", ph.GetAllPosts)
}

func Setup() *fiber.App {
	app := fiber.New()

	SetupRoutes(app)
	return app
}
