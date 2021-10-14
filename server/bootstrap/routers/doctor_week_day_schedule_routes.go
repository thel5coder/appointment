package routers

import (
	"github.com/gofiber/fiber/v2"
	"profira-backend/server/handlers"
	"profira-backend/server/middleware"
)

type DoctorWeekDayScheduleRoutes struct {
	RouteGroup fiber.Router
	Handler    handlers.Handler
}

func (route DoctorWeekDayScheduleRoutes) RegisterRoute() {
	handler := handlers.DoctorWeekDayScheduleHandler{Handler: route.Handler}
	jwtMiddleware := middleware.JwtMiddleware{UcContract: route.Handler.UseCaseContract}
	isActiveStatusMiddleware := middleware.VerifyIsActiveStatus{UcContract: jwtMiddleware.UcContract}

	doctorWeekDayScheduleRoutes := route.RouteGroup.Group("/schedule")
	doctorWeekDayScheduleRoutes.Use(jwtMiddleware.New)
	doctorWeekDayScheduleRoutes.Use(isActiveStatusMiddleware.IsActiveStatusMiddleware)
	doctorWeekDayScheduleRoutes.Get("", handler.BrowseBy)
	doctorWeekDayScheduleRoutes.Get("/doctor", handler.BrowseBy)
	doctorWeekDayScheduleRoutes.Post("", handler.Store)
}
