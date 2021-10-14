package handlers

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	"profira-backend/usecase"
)

type FileHandler struct {
	Handler
}

func (handler FileHandler) Add(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err)
	}

	uc := usecase.FileUseCase{UcContract: handler.UseCaseContract}
	res, err := uc.Add(file)

	return handler.SendResponse(ctx, res, nil, err)
}

func (handler FileHandler) Read(ctx *fiber.Ctx) error {
	ID := ctx.Params("id")

	uc := usecase.FileUseCase{UcContract: handler.UseCaseContract}
	res, err := uc.Read(ID)

	return handler.SendResponse(ctx, res, nil, err)
}
