package DtoResponse

import (
	"net/http"
	"profira-backend/usecase/viewmodel"
)

func SingleErrorMesssage(message string) []string {
	var errMessages []string

	errMessages = append(errMessages, message)

	return errMessages
}

func ErrorResponse(statusCode int, message interface{}) viewmodel.ResponseVm {
	responseVm := viewmodel.ResponseVm{
		Body: viewmodel.RespBodyVm{
			Message:    message,
			DataVm:     nil,
			Pagination: nil,
		},
		StatusCode: statusCode,
	}

	return responseVm
}

func SuccessResponse(dataVm interface{}, paginationVm interface{}) viewmodel.ResponseVm {
	responseVm := viewmodel.ResponseVm{
		Body: viewmodel.RespBodyVm{
			Message:    nil,
			DataVm:     dataVm,
			Pagination: paginationVm,
		},
		StatusCode: http.StatusOK,
	}

	return responseVm
}
