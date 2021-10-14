package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"profira-backend/server/requests"
	"profira-backend/usecase"
	"strconv"
)

type RoleHandler struct {
	Handler
}

func (handler RoleHandler) Browse(ctx *fiber.Ctx) error {
	search := ctx.Query("search")
	order := ctx.Query("order")
	sort := ctx.Query("sort")
	limit, _ := strconv.Atoi(ctx.Query("limit"))
	page, _ := strconv.Atoi(ctx.Query("page"))

	uc := usecase.RoleUseCase{UcContract: handler.UseCaseContract}
	res, pagination, err := uc.Browse(search, order, sort, page, limit)

	return handler.SendResponse(ctx, res, pagination, err)
}

func (handler RoleHandler) Read(ctx *fiber.Ctx) error {
	ID := ctx.Params("id")

	uc := usecase.RoleUseCase{UcContract: handler.UseCaseContract}
	res, err := uc.ReadBy("id",ID)

	return handler.SendResponse(ctx, res, nil, err)
}

func (handler RoleHandler) Edit(ctx *fiber.Ctx) error {
	ID := ctx.Params("id")
	input := new(requests.RoleRequest)

	if err := ctx.BodyParser(input); err != nil {
		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
	}
	if err := handler.Validate.Struct(input); err != nil {
		return handler.SendResponseErrorValidation(ctx, err.(validator.ValidationErrors))
	}

	uc := usecase.RoleUseCase{UcContract: handler.UseCaseContract}
	err := uc.Edit(ID, input)

	return handler.SendResponse(ctx, nil, nil, err)
}

func (handler RoleHandler) Add(ctx *fiber.Ctx) error {
	input := new(requests.RoleRequest)

	if err := ctx.BodyParser(input); err != nil {
		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
	}
	if err := handler.Validate.Struct(input); err != nil {
		return handler.SendResponseErrorValidation(ctx, err.(validator.ValidationErrors))
	}

	uc := usecase.RoleUseCase{UcContract: handler.UseCaseContract}
	err := uc.Add(input)

	return handler.SendResponse(ctx, nil, nil, err)
}

func (handler RoleHandler) Delete(ctx *fiber.Ctx) error {
	ID := ctx.Params("id")

	uc := usecase.RoleUseCase{UcContract: handler.UseCaseContract}
	err := uc.Delete(ID)

	return handler.SendResponse(ctx, nil, nil, err)
}