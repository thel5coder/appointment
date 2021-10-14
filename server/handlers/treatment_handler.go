package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"profira-backend/server/requests"
	"profira-backend/usecase"
	"strconv"
)

type TreatmentHandler struct {
	Handler
}

//browse
func (handler TreatmentHandler) Browse(ctx *fiber.Ctx) error {
	search := ctx.Query("search")
	order := ctx.Query("order")
	sort := ctx.Query("sort")
	limit, _ := strconv.Atoi(ctx.Query("limit"))
	page, _ := strconv.Atoi(ctx.Query("page"))

	uc := usecase.TreatmentUseCase{UcContract: handler.UseCaseContract}
	res, pagination, err := uc.Browse(search, order, sort, page, limit)

	return handler.SendResponse(ctx, res, pagination, err)
}

//browse all
func (handler TreatmentHandler) BrowseAll(ctx *fiber.Ctx) error {
	search := ctx.Query("search")

	uc := usecase.TreatmentUseCase{UcContract: handler.UseCaseContract}
	res, err := uc.BrowseAll(search)

	return handler.SendResponse(ctx, res, nil, err)
}

//browse by category
func (handler TreatmentHandler) BrowseByCategory(ctx *fiber.Ctx) error {
	categoryID := ctx.Params("id")

	uc := usecase.TreatmentUseCase{UcContract: handler.UseCaseContract}
	res, err := uc.BrowseByCategory(categoryID)

	return handler.SendResponse(ctx, res, nil, err)
}

//read by
func (handler TreatmentHandler) ReadByPk(ctx *fiber.Ctx) error {
	ID := ctx.Params("id")

	uc := usecase.TreatmentUseCase{UcContract: handler.UseCaseContract}
	res, err := uc.ReadBy("t.id", ID, "=")

	return handler.SendResponse(ctx, res, nil, err)
}

//edit
func (handler TreatmentHandler) Edit(ctx *fiber.Ctx) (err error) {
	ID := ctx.Params("id")
	input := new(requests.TreatmentRequest)

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
	uc := usecase.TreatmentUseCase{UcContract: handler.UseCaseContract}
	err = uc.Edit(input, ID)
	if err != nil {
		handler.UseCaseContract.TX.Rollback()

		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
	}
	handler.UseCaseContract.TX.Commit()

	return handler.SendResponse(ctx, nil, nil, err)
}

//add
func (handler TreatmentHandler) Add(ctx *fiber.Ctx) (err error) {
	input := new(requests.TreatmentRequest)

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
	uc := usecase.TreatmentUseCase{UcContract: handler.UseCaseContract}
	err = uc.Add(input)
	if err != nil {
		handler.UseCaseContract.TX.Rollback()

		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
	}
	handler.UseCaseContract.TX.Commit()

	return handler.SendResponse(ctx, nil, nil, err)
}

//delete
func (handler TreatmentHandler) Delete(ctx *fiber.Ctx) (err error) {
	ID := ctx.Params("id")

	handler.UseCaseContract.TX, err = handler.Db.Begin()
	if err != nil {
		handler.UseCaseContract.TX.Rollback()

		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
	}
	uc := usecase.TreatmentUseCase{UcContract: handler.UseCaseContract}
	err = uc.Delete(ID)
	if err != nil {
		handler.UseCaseContract.TX.Rollback()

		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
	}
	handler.UseCaseContract.TX.Commit()

	return handler.SendResponse(ctx, nil, nil, err)
}
