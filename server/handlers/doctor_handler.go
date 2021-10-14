package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"profira-backend/server/requests"
	"profira-backend/usecase"
	"strconv"
)

type DoctorHandler struct {
	Handler
}

//browse
func (handler DoctorHandler) Browse(ctx *fiber.Ctx) error {
	search := ctx.Query("search")
	order := ctx.Query("order")
	sort := ctx.Query("sort")
	limit, _ := strconv.Atoi(ctx.Query("limit"))
	page, _ := strconv.Atoi(ctx.Query("page"))

	uc := usecase.DoctorUseCase{UcContract: handler.UseCaseContract}
	res, pagination, err := uc.Browse(defaultDoctorRoleID, search, sort, order, page, limit)

	return handler.SendResponse(ctx, res, pagination, err)
}

//browse all
func (handler DoctorHandler) BrowseAll(ctx *fiber.Ctx) error {
	search := ctx.Query("search")

	uc := usecase.DoctorUseCase{UcContract: handler.UseCaseContract}
	res, err := uc.BrowseAll(defaultDoctorRoleID, search)

	return handler.SendResponse(ctx, res, nil, err)
}

//read
func (handler DoctorHandler) ReadByPk(ctx *fiber.Ctx) error {
	ID := ctx.Params("id")

	uc := usecase.DoctorUseCase{UcContract: handler.UseCaseContract}
	res, err := uc.ReadBy("s.id", ID, "=")

	return handler.SendResponse(ctx, res, nil, err)
}

//edit
func (handler DoctorHandler) Edit(ctx *fiber.Ctx) (err error) {
	input := new(requests.StaffRequest)
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
	uc := usecase.StaffUseCase{UcContract: handler.UseCaseContract}
	err = uc.Edit(input, ID)
	if err != nil {
		handler.UseCaseContract.TX.Rollback()

		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
	}
	handler.UseCaseContract.TX.Commit()

	return handler.SendResponse(ctx, nil, nil, err)
}

//add
func (handler DoctorHandler) Add(ctx *fiber.Ctx) (err error) {
	input := new(requests.StaffRequest)

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
	uc := usecase.StaffUseCase{UcContract: handler.UseCaseContract}
	err = uc.Add(input, defaultDoctorRoleID)
	if err != nil {
		handler.UseCaseContract.TX.Rollback()

		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
	}
	handler.UseCaseContract.TX.Commit()

	return handler.SendResponse(ctx, nil, nil, err)
}

//delete
func (handler DoctorHandler) Delete(ctx *fiber.Ctx) (err error) {
	ID := ctx.Params("id")

	handler.UseCaseContract.TX, err = handler.Db.Begin()
	if err != nil {
		handler.UseCaseContract.TX.Rollback()

		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
	}
	uc := usecase.StaffUseCase{UcContract: handler.UseCaseContract}
	err = uc.DeleteBy("s.id", ID, "=")
	if err != nil {
		handler.UseCaseContract.TX.Rollback()

		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
	}
	handler.UseCaseContract.TX.Commit()

	return handler.SendResponse(ctx, nil, nil, err)
}
