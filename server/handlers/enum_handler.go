package handlers

import (
	"github.com/gofiber/fiber/v2"
	"profira-backend/usecase"
)

type EnumHandler struct {
	Handler
}

func (handler EnumHandler) GetSexEnums(ctx *fiber.Ctx) error {
	uc := usecase.EnumUseCase{UcContract: handler.UseCaseContract}
	res := uc.GetSexEnum()

	return handler.SendResponse(ctx, res, nil, nil)
}
