package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"profira-backend/server/requests"
	"profira-backend/usecase"
	"strconv"
)

type ProcedureHandler struct {
	Handler
}

func (handler ProcedureHandler) Browse(ctx *fiber.Ctx) error{
	search := ctx.Query("search")
	order := ctx.Query("order")
	sort := ctx.Query("sort")
	limit, _ := strconv.Atoi(ctx.Query("limit"))
	page, _ := strconv.Atoi(ctx.Query("page"))

	uc := usecase.ProcedureUseCase{UcContract: handler.UseCaseContract}
	res, pagination, err := uc.Browse(search, order, sort, page, limit)

	return handler.SendResponse(ctx, res, pagination, err)
}

func (handler ProcedureHandler) BrowseAll(ctx *fiber.Ctx) error{
	search := ctx.Query("search")

	uc := usecase.ProcedureUseCase{UcContract: handler.UseCaseContract}
	res, err := uc.BrowseAll(search)

	return handler.SendResponse(ctx, res, nil, err)
}

func (handler ProcedureHandler) ReadByPk(ctx *fiber.Ctx) error{
	ID := ctx.Params("id")

	uc := usecase.ProcedureUseCase{UcContract: handler.UseCaseContract}
	res, err := uc.ReadBy("id",ID,"=")

	return handler.SendResponse(ctx, res, nil, err)
}

func (handler ProcedureHandler) Edit(ctx *fiber.Ctx) error{
	ID := ctx.Params("id")
	input := new(requests.ProcedureRequest)

	if err := ctx.BodyParser(input); err != nil {
		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
	}
	if err := handler.Validate.Struct(input); err != nil {
		return handler.SendResponseErrorValidation(ctx, err.(validator.ValidationErrors))
	}

	uc := usecase.ProcedureUseCase{UcContract: handler.UseCaseContract}
	err := uc.Edit(input, ID)

	return handler.SendResponse(ctx, nil, nil, err)
}

func (handler ProcedureHandler) Add(ctx *fiber.Ctx) error{
	input := new(requests.ProcedureRequest)

	if err := ctx.BodyParser(input); err != nil {
		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
	}
	if err := handler.Validate.Struct(input); err != nil {
		return handler.SendResponseErrorValidation(ctx, err.(validator.ValidationErrors))
	}

	uc := usecase.ProcedureUseCase{UcContract: handler.UseCaseContract}
	err := uc.Add(input)

	return handler.SendResponse(ctx, nil, nil, err)
}

func (handler ProcedureHandler) DeleteByPk(ctx *fiber.Ctx) error{
	ID := ctx.Params("id")

	uc := usecase.ProcedureUseCase{UcContract: handler.UseCaseContract}
	err := uc.DeleteBy("id",ID,"=")

	return handler.SendResponse(ctx, nil, nil, err)
}
