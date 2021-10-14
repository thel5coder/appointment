package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"profira-backend/server/requests"
	"profira-backend/usecase"
)

type CategoryHandler struct {
	Handler
}

//browse tree view
func(handler CategoryHandler) BrowseTreeView(ctx *fiber.Ctx) error{
	uc := usecase.CategoryUseCase{UcContract: handler.UseCaseContract}

	res, err := uc.BrowseTreeView()

	return handler.SendResponse(ctx, res, nil, err)
}

//browse by parent
func (handler CategoryHandler) BrowseByParent(ctx *fiber.Ctx) error {
	parentID := ctx.Query("parent_id")
	uc := usecase.CategoryUseCase{UcContract: handler.UseCaseContract}

	res, err := uc.BrowseByParent(parentID)

	return handler.SendResponse(ctx, res, nil, err)
}

//browse category for mobile
func (handler CategoryHandler) BrowseShortcutCategory(ctx *fiber.Ctx) error {
	uc := usecase.CategoryUseCase{UcContract: handler.UseCaseContract}

	res, err := uc.BrowseShortcutCategory()

	return handler.SendResponse(ctx, res, nil, err)
}

//read
func (handler CategoryHandler) ReadByID(ctx *fiber.Ctx) error {
	ID := ctx.Params("id")

	uc := usecase.CategoryUseCase{UcContract: handler.UseCaseContract}
	res, err := uc.ReadBy("c.id", ID, "=")

	return handler.SendResponse(ctx, res, nil, err)
}

//edit category
func (handler CategoryHandler) Edit(ctx *fiber.Ctx) (err error) {
	input := new(requests.CategoryRequest)
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
	uc := usecase.CategoryUseCase{UcContract: handler.UseCaseContract}
	err = uc.Edit(input, ID)
	if err != nil {
		handler.UseCaseContract.TX.Rollback()

		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
	}
	handler.UseCaseContract.TX.Commit()

	return handler.SendResponse(ctx, nil, nil, err)
}

//add
func (handler CategoryHandler) Add(ctx *fiber.Ctx) (err error) {
	input := new(requests.CategoryRequest)

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
	uc := usecase.CategoryUseCase{UcContract: handler.UseCaseContract}
	err = uc.Add(input)
	if err != nil {
		handler.UseCaseContract.TX.Rollback()

		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
	}
	handler.UseCaseContract.TX.Commit()

	return handler.SendResponse(ctx, nil, nil, err)
}

func (handler CategoryHandler) Delete(ctx *fiber.Ctx) error {
	ID := ctx.Params("id")

	uc := usecase.CategoryUseCase{UcContract: handler.UseCaseContract}
	err := uc.Delete(ID)

	return handler.SendResponse(ctx, nil, nil, err)
}
