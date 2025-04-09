package router

import (
	"cache_engine_httpserver/internal/api/config"
	"cache_engine_httpserver/internal/api/http"
	"cache_engine_httpserver/internal/api/model"

	"github.com/gofiber/fiber/v3"
)

func HandleRoute(app *fiber.App, ctx *model.CacheAppContext) {
	app.Get(config.BASE_URL_NAME+"/get", func(c fiber.Ctx) error {
		return http.GetCache(c, ctx)
	})

	app.Post(config.BASE_URL_NAME+"/create", func(c fiber.Ctx) error {
		return http.CreateCache(c, ctx)
	})

	app.Delete(config.BASE_URL_NAME+"/delete/:key", func(c fiber.Ctx) error {
		return http.DeleteCache(c, ctx)
	})

	app.Get(config.BASE_URL_NAME+"/exists/:key", func(c fiber.Ctx) error {
		return http.IsCacheExists(c, ctx)
	})
}
