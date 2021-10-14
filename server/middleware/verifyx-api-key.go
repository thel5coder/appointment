package middleware

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"profira-backend/helpers/messages"
	"profira-backend/server/handlers"
	"profira-backend/usecase"
	"time"
)

type VerifyXAPIKey struct {
	*usecase.UcContract
}

func (verifyXAPIKey VerifyXAPIKey) XAPIKeyVerify(ctx *fiber.Ctx) error {
	apiHandler := handlers.Handler{UseCaseContract: verifyXAPIKey.UcContract}
	var sha1 = sha1.New()
	now := time.Now().UTC().Format("2006-01-02T15:04Z07:00")

	sha1.Write([]byte(now))
	sha1EncryptedStr := sha1.Sum(nil)
	encrypted := fmt.Sprintf("%x", sha1EncryptedStr)
	fmt.Println(encrypted)
	xAPIKey := ctx.Get("x-api-key")

	if encrypted != xAPIKey {
		return apiHandler.SendResponseUnauthorized(ctx, errors.New(messages.Unauthorized))
	}

	return ctx.Next()
}
