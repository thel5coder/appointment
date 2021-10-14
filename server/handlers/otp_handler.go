package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"profira-backend/server/requests"
	"profira-backend/usecase"
)

type OTPHandler struct {
	Handler
}

func (handler OTPHandler) RequestOTP(ctx *fiber.Ctx) error {
	input := new(requests.OTPRequest)

	if err := ctx.BodyParser(input); err != nil {
		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
	}
	if err := handler.Validate.Struct(input); err != nil {
		return handler.SendResponseErrorValidation(ctx, err.(validator.ValidationErrors))
	}

	uc := usecase.OtpUseCase{UcContract: handler.UseCaseContract}
	err := uc.RequestOTP(input.MobilePhone)

	return handler.SendResponse(ctx, nil, nil, err)
}

func (handler OTPHandler) SubmitOTP(ctx *fiber.Ctx) error{
	return nil
}

func(handler OTPHandler) RequestXAPIKey(ctx *fiber.Ctx) error {
	uc := usecase.OtpUseCase{UcContract: handler.UseCaseContract}
	res := uc.RequestXAPIKey()

	return handler.SendResponse(ctx, res, nil, nil)
}
