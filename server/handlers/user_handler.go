package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"profira-backend/server/requests"
	"profira-backend/usecase"
	"strconv"
)

type UserHandler struct {
	Handler
}

func (handler UserHandler) Browse(ctx *fiber.Ctx) error {
	search := ctx.Query("search")
	order := ctx.Query("order")
	sort := ctx.Query("sort")
	limit, _ := strconv.Atoi(ctx.Query("limit"))
	page, _ := strconv.Atoi(ctx.Query("page"))

	uc := usecase.UserUseCase{UcContract: handler.UseCaseContract}
	res, pagination, err := uc.Browse(search, order, sort, page, limit)

	return handler.SendResponse(ctx, res, pagination, err)
}

func (handler UserHandler) BrowseAll(ctx *fiber.Ctx) error {
	uc := usecase.UserUseCase{UcContract: handler.UseCaseContract}
	res, err := uc.BrowseAllBy("", "")

	return handler.SendResponse(ctx, res, nil, err)
}

func (handler UserHandler) ReadByPk(ctx *fiber.Ctx) error {
	ID := ctx.Params("id")

	uc := usecase.UserUseCase{UcContract: handler.UseCaseContract}
	res, err := uc.ReadBy("u.id", ID, "=")

	return handler.SendResponse(ctx, res, nil, err)
}

func (handler UserHandler) Edit(ctx *fiber.Ctx) (err error) {
	ID := ctx.Params("id")
	input := new(requests.UserRequest)

	if err := ctx.BodyParser(input); err != nil {
		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
	}
	if err := handler.Validate.Struct(input); err != nil {
		return handler.SendResponseErrorValidation(ctx, err.(validator.ValidationErrors))
	}

	handler.UseCaseContract.TX, err = handler.Db.Begin()
	if err != nil {
		handler.UseCaseContract.TX.Rollback()

		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
	}
	uc := usecase.UserUseCase{UcContract: handler.UseCaseContract}
	err = uc.Edit(ID, input)
	if err != nil {
		handler.UseCaseContract.TX.Rollback()

		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
	}
	handler.UseCaseContract.TX.Commit()

	return handler.SendResponse(ctx, nil, nil, err)
}

func (handler UserHandler) Add(ctx *fiber.Ctx) (err error) {
	input := new(requests.UserRequest)

	if err := ctx.BodyParser(input); err != nil {
		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
	}
	if err := handler.Validate.Struct(input); err != nil {
		return handler.SendResponseErrorValidation(ctx, err.(validator.ValidationErrors))
	}

	handler.UseCaseContract.TX, err = handler.Db.Begin()
	if err != nil {
		handler.UseCaseContract.TX.Rollback()

		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
	}
	uc := usecase.UserUseCase{UcContract: handler.UseCaseContract}
	_, err = uc.Add(input)
	if err != nil {
		handler.UseCaseContract.TX.Rollback()

		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
	}
	handler.UseCaseContract.TX.Commit()

	return handler.SendResponse(ctx, nil, nil, err)
}

func (handler UserHandler) DeleteByPk(ctx *fiber.Ctx) (err error) {
	ID := ctx.Params("id")

	handler.UseCaseContract.TX, err = handler.Db.Begin()
	if err != nil {
		handler.UseCaseContract.TX.Rollback()

		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err)
	}
	uc := usecase.UserUseCase{UcContract: handler.UseCaseContract}
	err = uc.DeleteBy("id", ID, "=")
	if err != nil {
		handler.UseCaseContract.TX.Rollback()

		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err)
	}
	handler.UseCaseContract.TX.Commit()

	return handler.SendResponse(ctx, nil, nil, err)
}

func (handler UserHandler) ReadProfile(ctx *fiber.Ctx) error {
	uc := usecase.UserUseCase{UcContract: handler.UseCaseContract}
	res, err := uc.ReadProfile()

	return handler.SendResponse(ctx, res, nil, err)
}
