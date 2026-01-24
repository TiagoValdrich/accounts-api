package config

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rs/zerolog/log"
	"github.com/tiagovaldrich/accounts-api/internal/pkg/cerror"
)

type AppRouter struct {
	app *fiber.App
}

func NewRouter() *AppRouter {
	router := &AppRouter{}

	app := fiber.New(fiber.Config{
		AppName:               AppName,
		DisableStartupMessage: true,
		ErrorHandler:          router.ErrorHandler,
	})

	app.Use(cors.New())
	app.Get("/status", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	router.app = app

	return router
}

func (router *AppRouter) ErrorHandler(c *fiber.Ctx, err error) error {
	log.Err(err).Send()

	var customError *cerror.Error
	if errors.As(err, &customError) {
		if customError.Status == 0 {
			customError.Status = fiber.StatusInternalServerError
		}

		if customError.Message == "" {
			customError.Message = "Internal server error"
		}

		return c.Status(customError.Status).JSON(customError)
	}

	return c.SendStatus(fiber.StatusInternalServerError)
}

func (router *AppRouter) Start() error {
	log.Info().Msg("starting server on port " + AppPort)

	if err := router.app.Listen(AppPort); err != nil {
		panic(err)
	}

	return nil
}

func (router *AppRouter) GetApp() *fiber.App {
	return router.app
}
