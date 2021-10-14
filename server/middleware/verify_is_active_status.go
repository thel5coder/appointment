package middleware

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"profira-backend/helpers/messages"
	"profira-backend/server/handlers"
	"profira-backend/usecase"
)

type VerifyIsActiveStatus struct {
	*usecase.UcContract
}

func(verifyIsActiveStatus VerifyIsActiveStatus) IsActiveStatusMiddleware(ctx *fiber.Ctx) error{
	apiHandler := handlers.Handler{UseCaseContract: verifyIsActiveStatus.UcContract}
	userUc := usecase.UserUseCase{UcContract:verifyIsActiveStatus.UcContract}
	user,err := userUc.ReadBy("u.id",userUc.UserID,"=")
	if err != nil {
		return apiHandler.SendResponseUnauthorized(ctx, errors.New(messages.Unauthorized))
	}
	if !user.IsActive {
		return apiHandler.SendResponseUnauthorized(ctx, errors.New(messages.Unauthorized))
	}

	return ctx.Next()
}
