package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"profira-backend/server/requests"
	"profira-backend/usecase"
	"strconv"
)

type ClinicHandler struct {
	Handler
}

func (handler ClinicHandler) Browse(ctx *fiber.Ctx) error {
	search := ctx.Query("search")
	order := ctx.Query("order")
	sort := ctx.Query("sort")
	limit, _ := strconv.Atoi(ctx.Query("limit"))
	page, _ := strconv.Atoi(ctx.Query("page"))

	uc := usecase.ClinicUseCase{UcContract: handler.UseCaseContract}
	res, pagination, err := uc.Browse(search, order, sort, page, limit)

	return handler.SendResponse(ctx, res, pagination, err)
}

func (handler ClinicHandler) BrowseAll(ctx *fiber.Ctx) error {
	search := ctx.Query("search")

	uc := usecase.ClinicUseCase{UcContract: handler.UseCaseContract}
	res, err := uc.BrowseAll(search)

	return handler.SendResponse(ctx, res, nil, err)
}

func (handler ClinicHandler) ReadByPk(ctx *fiber.Ctx) error {
	ID := ctx.Params("id")

	uc := usecase.ClinicUseCase{UcContract: handler.UseCaseContract}
	res, err := uc.ReadBy("id", ID, "=")

	return handler.SendResponse(ctx, res, nil, err)
}

func (handler ClinicHandler) Edit(ctx *fiber.Ctx) error {
	input := new(requests.ClinicRequest)
	ID := ctx.Params("id")

	if err := ctx.BodyParser(input); err != nil {
		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
	}
	if err := handler.Validate.Struct(input); err != nil {
		return handler.SendResponseErrorValidation(ctx, err.(validator.ValidationErrors))
	}

	uc := usecase.ClinicUseCase{UcContract: handler.UseCaseContract}
	err := uc.Edit(input, ID)

	return handler.SendResponse(ctx, nil, nil, err)
}

func (handler ClinicHandler) Add(ctx *fiber.Ctx) error {
	input := new(requests.ClinicRequest)

	if err := ctx.BodyParser(input); err != nil {
		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
	}
	if err := handler.Validate.Struct(input); err != nil {
		return handler.SendResponseErrorValidation(ctx, err.(validator.ValidationErrors))
	}

	uc := usecase.ClinicUseCase{UcContract: handler.UseCaseContract}
	err := uc.Add(input)

	return handler.SendResponse(ctx, nil, nil, err)
}

func (handler ClinicHandler) Delete(ctx *fiber.Ctx) error {
	ID := ctx.Params("id")

	uc := usecase.ClinicUseCase{UcContract: handler.UseCaseContract}
	err := uc.Delete(ID)

	return handler.SendResponse(ctx, nil, nil, err)
}
