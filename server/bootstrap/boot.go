package bootstrap

import (
	"database/sql"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtFiber "github.com/gofiber/jwt/v2"
	"profira-backend/helpers/jwe"
	"profira-backend/helpers/jwt"
	"profira-backend/usecase"
)

type Bootstrap struct {
	App             *fiber.App
	Db              *sql.DB
	UseCaseContract usecase.UcContract
	Jwe             jwe.Credential
	Validator       *validator.Validate
	Translator      ut.Translator
	JwtConfig       jwtFiber.Config
	JwtCred         jwt.JwtCredential
}
