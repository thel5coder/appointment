package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"profira-backend/server/requests"
	"profira-backend/usecase"
)

type DoctorWeekDayScheduleHandler struct {
	Handler
}

//browse by
func (handler DoctorWeekDayScheduleHandler) BrowseBy(ctx *fiber.Ctx) error {
	filters := make(map[string]interface{})
	if ctx.Query("clinic_id") != "" {
		filters["clinic_id"] = ctx.Query("clinic_id")
	}
	if ctx.Query("doctor_id") != "" {
		filters["doctor_id"] = ctx.Query("doctor_id")
	}

	uc := usecase.DoctorScheduleWeekDayUseCase{UcContract: handler.UseCaseContract}
	res, err := uc.BrowseBy(filters)

	return handler.SendResponse(ctx, res, nil, err)
}

//store
func (handler DoctorWeekDayScheduleHandler) Store(ctx *fiber.Ctx) (err error) {
	input := new(requests.DoctorScheduleRequest)

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
	uc := usecase.DoctorScheduleWeekDayUseCase{UcContract: handler.UseCaseContract}
	err = uc.Store(input)
	if err != nil {
		handler.UseCaseContract.TX.Rollback()

		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
	}
	handler.UseCaseContract.TX.Commit()

	return handler.SendResponse(ctx, nil, nil, err)
}


