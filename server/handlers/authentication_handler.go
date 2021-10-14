package handlers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"profira-backend/server/requests"
	"profira-backend/usecase"
)

type AuthenticationHandler struct {
	Handler
}

func (handler AuthenticationHandler) Login(ctx *fiber.Ctx) error {
	input := new(requests.LoginRequest)

	if err := ctx.BodyParser(input); err != nil {
		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
	}
	fmt.Println(input)
	if err := handler.Validate.Struct(input); err != nil {
		return handler.SendResponseErrorValidation(ctx, err.(validator.ValidationErrors))
	}

	uc := usecase.AuthenticationUseCase{UcContract: handler.UseCaseContract}
	res, err := uc.Login(input)
	if err != nil {
		return handler.SendResponseUnauthorized(ctx, err)
	}

	return handler.SendResponse(ctx, res, nil, err)
}

func (handler AuthenticationHandler) Registration(ctx *fiber.Ctx) (err error) {
	input := new(requests.RegisterRequest)

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

	uc := usecase.AuthenticationUseCase{UcContract: handler.UseCaseContract}
	res, err := uc.Registration(input)
	if err != nil {
		handler.UseCaseContract.TX.Rollback()

		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
	}
	handler.UseCaseContract.TX.Commit()

	return handler.SendResponse(ctx, res, nil, err)
}

func (handler AuthenticationHandler) ActivationCustomer(ctx *fiber.Ctx) error {
	input := new(requests.ActivationCustomerRequest)

	if err := ctx.BodyParser(input); err != nil {
		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
	}
	if err := handler.Validate.Struct(input); err != nil {
		return handler.SendResponseErrorValidation(ctx, err.(validator.ValidationErrors))
	}

	uc := usecase.AuthenticationUseCase{UcContract: handler.UseCaseContract}
	err := uc.ActivationCustomer(input)

	return handler.SendResponse(ctx, nil, nil, err)
}

func (handler AuthenticationHandler) ForgotPassword(ctx *fiber.Ctx) error {
	input := new(requests.ForgotPasswordRequest)

	if err := ctx.BodyParser(input); err != nil {
		return handler.SendResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
	}
	if err := handler.Validate.Struct(input); err != nil {
		return handler.SendResponseErrorValidation(ctx, err.(validator.ValidationErrors))
	}

	uc := usecase.AuthenticationUseCase{UcContract: handler.UseCaseContract}
	err := uc.ForgotPassword(input.Email)

	return handler.SendResponse(ctx, nil, nil, err)
}
