package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"profira-backend/server/requests"
	"profira-backend/usecase"
	"strconv"
)

type PromotionHandler struct {
	Handler
}

//browse
func (handler PromotionHandler) Browse(ctx *fiber.Ctx) error {
	search := ctx.Query("search")
	order := ctx.Query("order")
	sort := ctx.Query("sort")
	limit, _ := strconv.Atoi(ctx.Query("limit"))
	page, _ := strconv.Atoi(ctx.Query("page"))

	uc := usecase.PromotionUseCase{UcContract: handler.UseCaseContract}
	res, pagination, err := uc.Browse(search, order, sort, page, limit)

	return handler.SendResponse(ctx, res, pagination, err)
}

//browse active promotion
func (handler PromotionHandler) BrowseActivePromotion(ctx *fiber.Ctx) error {
	uc := usecase.PromotionUseCase{UcContract: handler.UseCaseContract}
	res, err := uc.BrowseActivePromotion()

	return handler.SendResponse(ctx, res, nil, err)
}

//read by Id
func (handler PromotionHandler) ReadBy(ctx *fiber.Ctx) error {
	ID := ctx.Params("id")

	uc := usecase.PromotionUseCase{UcContract: handler.UseCaseContract}
	res, err := uc.ReadBy("p.id", ID, "=")

	return handler.SendResponse(ctx, res, nil, err)
}

//edit
func (handler PromotionHandler) Edit(ctx *fiber.Ctx) (err error) {
	input := new(requests.PromotionRequest)
	ID := ctx.Params("id")

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
	uc := usecase.PromotionUseCase{UcContract: handler.UseCaseContract}
	err = uc.Edit(input, ID)
	if err != nil {
		handler.UseCaseContract.TX.Rollback()

		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
	}
	handler.UseCaseContract.TX.Commit()

	return handler.SendResponse(ctx, nil, nil, err)
}

//add
func (handler PromotionHandler) Add(ctx *fiber.Ctx) (err error) {
	input := new(requests.PromotionRequest)

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
	uc := usecase.PromotionUseCase{UcContract: handler.UseCaseContract}
	err = uc.Add(input)
	if err != nil {
		handler.UseCaseContract.TX.Rollback()

		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
	}
	handler.UseCaseContract.TX.Commit()

	return handler.SendResponse(ctx, nil, nil, err)
}

//delete
func (handler PromotionHandler) Delete(ctx *fiber.Ctx) (err error) {
	ID := ctx.Params("id")

	handler.UseCaseContract.TX, err = handler.Db.Begin()
	if err != nil {
		handler.UseCaseContract.TX.Rollback()

		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
	}
	uc := usecase.PromotionUseCase{UcContract: handler.UseCaseContract}
	err = uc.Delete(ID)
	if err != nil {
		handler.UseCaseContract.TX.Rollback()

		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
	}
	handler.UseCaseContract.TX.Commit()

	return handler.SendResponse(ctx, nil, nil, err)
}
