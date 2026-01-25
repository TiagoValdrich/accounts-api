package accounts

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/tiagovaldrich/accounts-api/internal/pkg/cerror"
)

type httpHandler struct {
	service Servicer
}

func NewHTTPHandler(app *fiber.App, service Servicer) {
	httpHandler := &httpHandler{
		service: service,
	}

	routeGroup := app.Group("/accounts")
	routeGroup.Post("/", httpHandler.createAccount)
}

func (h *httpHandler) createAccount(c *fiber.Ctx) error {
	var body CreateAccountRequest

	if err := c.BodyParser(&body); err != nil {
		return cerror.New(cerror.Params{
			Status:  http.StatusBadRequest,
			Message: "Invalid account payload",
		})
	}

	createAccountResult, err := h.service.CreateAccount(c.Context(), body)
	if err != nil {
		log.Err(err).Msg("failed to create account")

		return err
	}

	return c.Status(http.StatusOK).JSON(toAccountCreatedResponse(createAccountResult))
}
