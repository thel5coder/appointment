package handlers

import (
	"database/sql"
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jwtFiber "github.com/gofiber/jwt/v2"
	"net/http"
	"profira-backend/helpers/DtoResponse"
	"profira-backend/helpers/jwe"
	"profira-backend/helpers/str"
	"profira-backend/usecase"
	"profira-backend/usecase/viewmodel"
	"strings"
)

type Handler struct {
	App             *fiber.App
	UseCaseContract *usecase.UcContract
	Jwe             jwe.Credential
	Db              *sql.DB
	Validate        *validator.Validate
	Translator      ut.Translator
	JwtConfig        jwtFiber.Config
}

const (
	defaultDoctorRoleID = "d3ab093c-6866-46d3-95dc-67dd7ac3ffcb"
	defaultAdminRoleID  = "ecfee405-0da0-4d73-b9fc-3e42c9bb3e00"
	defaultNurseRoleID  = "9eaad38d-130d-49e2-b687-61e6ce77adc3"
)

func (h Handler) SendResponse(ctx *fiber.Ctx, data interface{}, pagination interface{}, err error) error {
	response := DtoResponse.SuccessResponse(data, pagination)
	if err != nil {
		response = DtoResponse.ErrorResponse(http.StatusUnprocessableEntity, err.Error())
	}

	return ctx.Status(response.StatusCode).JSON(response.Body)
}

func (h Handler) SendResponseBadRequest(ctx *fiber.Ctx, statusCode int, err interface{}) error {
	response := DtoResponse.ErrorResponse(statusCode, err)

	return ctx.Status(response.StatusCode).JSON(response.Body)
}

func (h Handler) SendResponseErrorValidation(ctx *fiber.Ctx, error validator.ValidationErrors) error {
	errorMessages := h.ExtractErrorValidationMessages(error)

	return h.SendResponseBadRequest(ctx, http.StatusBadRequest, errorMessages)
}

func (h Handler) SendResponseUnauthorized(ctx *fiber.Ctx, err error) error {
	response := DtoResponse.ErrorResponse(http.StatusUnauthorized, err.Error())

	return ctx.Status(response.StatusCode).JSON(response.Body)
}

func (h Handler) ResponseBadRequest(error string) viewmodel.ResponseVm {
	responseVm := viewmodel.ResponseVm{
		Body: viewmodel.RespBodyVm{
			Message:    error,
			DataVm:     nil,
			Pagination: nil,
		},
		StatusCode: http.StatusBadRequest,
	}

	return responseVm
}

func (h Handler) ResponseValidationError(error validator.ValidationErrors) viewmodel.ResponseVm {
	errorMessage := map[string][]string{}
	errorTranslation := error.Translate(h.Translator)

	for _, err := range error {
		errKey := str.Underscore(err.StructField())
		errorMessage[errKey] = append(
			errorMessage[errKey],
			strings.Replace(errorTranslation[err.Namespace()], err.StructField(), err.StructField(), -1),
		)
	}

	response := viewmodel.ResponseVm{
		Body: viewmodel.RespBodyVm{
			Message:    errorMessage,
			DataVm:     nil,
			Pagination: nil,
		},
		StatusCode: http.StatusBadRequest,
	}

	return response
}

func (h Handler) ExtractErrorValidationMessages(error validator.ValidationErrors) map[string][]string {
	errorMessage := map[string][]string{}
	errorTranslation := error.Translate(h.Translator)

	for _, err := range error {
		fmt.Println(errorTranslation[err.Namespace()])
		fmt.Println(err.StructField())
		fmt.Println(strings.Replace(errorTranslation[err.Namespace()], err.StructField(), err.StructField(), -1))
		errKey := str.Underscore(err.StructField())
		errorMessage[errKey] = append(errorMessage[errKey],
			strings.Replace(errorTranslation[err.Namespace()], err.StructField(), err.StructField(), -1),
		)
	}

	return errorMessage
}
