package handlers

import (
	"github.com/gofiber/fiber/v2"
	"profira-backend/usecase"
)

type DoctorAvailabilityHandler struct {
	Handler
}

func (handler DoctorAvailabilityHandler) BrowseDoctorAvailability(ctx *fiber.Ctx) error {
	clinicId := ctx.Query("clinic_id")
	//date := ctx.Query("date")

	uc := usecase.DoctorUseCase{UcContract: handler.UseCaseContract}
	res, err := uc.BrowseAllByClinic(clinicId, defaultDoctorRoleID)

	return handler.SendResponse(ctx, res, nil, err)
}
