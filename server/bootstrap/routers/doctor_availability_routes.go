package routers

import (
	"github.com/gofiber/fiber/v2"
	"profira-backend/server/handlers"
	"profira-backend/server/middleware"
)

type DoctorAvailabilityRoutes struct {
	RouteGroup fiber.Router
	Handler    handlers.Handler
}

func (route DoctorAvailabilityRoutes) RegisterRouter() {
	handler := handlers.DoctorAvailabilityHandler{Handler: route.Handler}
	jwtMiddleware := middleware.JwtMiddleware{UcContract: route.Handler.UseCaseContract}
	isActiveStatusMiddleware := middleware.VerifyIsActiveStatus{UcContract: jwtMiddleware.UcContract}

	doctorAvailabilityRoutes := route.RouteGroup.Group("/schedule")
	doctorAvailabilityRoutes.Use(jwtMiddleware.New)
	doctorAvailabilityRoutes.Use(isActiveStatusMiddleware.IsActiveStatusMiddleware)
	doctorAvailabilityRoutes.Get("", handler.BrowseDoctorAvailability)
}
