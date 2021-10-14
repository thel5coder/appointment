package routers

import (
	"github.com/gofiber/fiber/v2"
	"profira-backend/server/handlers"
	"profira-backend/server/middleware"
)

type ClinicRoutes struct {
	RouteGroup fiber.Router
	Handler    handlers.Handler
}

func (route ClinicRoutes) RegisterRoute() {
	handler := handlers.ClinicHandler{Handler: route.Handler}
	jwtMiddleware := middleware.JwtMiddleware{UcContract: route.Handler.UseCaseContract}
	isActiveStatusMiddleware := middleware.VerifyIsActiveStatus{UcContract: jwtMiddleware.UcContract}

	clinicRoutes := route.RouteGroup.Group("/clinic")
	clinicRoutes.Use(jwtMiddleware.New)
	clinicRoutes.Use(isActiveStatusMiddleware.IsActiveStatusMiddleware)
	clinicRoutes.Get("", handler.Browse)
	clinicRoutes.Get("/all", handler.BrowseAll)
	clinicRoutes.Get("/:id", handler.ReadByPk)
	clinicRoutes.Put("/:id", handler.Edit)
	clinicRoutes.Post("", handler.Add)
	clinicRoutes.Delete("/:id", handler.Delete)
}
