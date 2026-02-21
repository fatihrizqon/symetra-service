package util

import (
	"github.com/fatihrizqon/symetra-service/internal/delivery/http/response"
	"github.com/gofiber/fiber/v2"
)

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func HandleError(ctx *fiber.Ctx, status int, err error) {
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(response.JSON{
			Status:  status,
			Message: err.Error(),
			Errors:  err,
		})
	}
}
