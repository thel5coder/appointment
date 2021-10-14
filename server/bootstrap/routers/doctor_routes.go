package routers

import (
	"github.com/gofiber/fiber/v2"
	"profira-backend/server/handlers"
	"profira-backend/server/middleware"
)

type DoctorRoutes struct {
	RouteGroup fiber.Router
	Handler    handlers.Handler
}

func (route DoctorRoutes) RegisterRoute() {
	handler := handlers.DoctorHandler{Handler: route.Handler}
	jwtMiddleware := middleware.JwtMiddleware{UcContract: route.Handler.UseCaseContract}
	isActiveStatusMiddleware := middleware.VerifyIsActiveStatus{UcContract: jwtMiddleware.UcContract}

	doctorRoutes := route.RouteGroup.Group("/doctor")
	doctorRoutes.Use(jwtMiddleware.New)
	doctorRoutes.Use(isActiveStatusMiddleware.IsActiveStatusMiddleware)
	doctorRoutes.Get("", handler.Browse)
	doctorRoutes.Get("/all", handler.BrowseAll)
	doctorRoutes.Get("/:id", handler.ReadByPk)
	doctorRoutes.Put("/:id", handler.Edit)
	doctorRoutes.Post("", handler.Add)
	doctorRoutes.Delete("/:id", handler.Delete)
}
