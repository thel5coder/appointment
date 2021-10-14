package routers

import (
	"github.com/gofiber/fiber/v2"
	"profira-backend/server/handlers"
	"profira-backend/server/middleware"
)

type OTPRoutes struct {
	RouteGroup fiber.Router
	Handler handlers.Handler
}

func (route OTPRoutes) RegisterRoute(){
	handler := handlers.OTPHandler{Handler:route.Handler}
	jwtMiddleware := middleware.JwtMiddleware{UcContract:handler.UseCaseContract}

	route.RouteGroup.Get("/x-api-key",handler.RequestXAPIKey)

	otpRoutes := route.RouteGroup.Group("/otp")
	otpRoutes.Use(jwtMiddleware.New)
	otpRoutes.Post("",handler.RequestOTP)
}
