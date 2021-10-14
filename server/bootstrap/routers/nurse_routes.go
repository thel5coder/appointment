package routers

import (
	"github.com/gofiber/fiber/v2"
	"profira-backend/server/handlers"
	"profira-backend/server/middleware"
)

type NurseRoutes struct {
	RouteGroup fiber.Router
	Handler    handlers.Handler
}

func (route NurseRoutes) RegisterRoute() {
	handler := handlers.NurseHandler{Handler: route.Handler}
	jwtMiddleware := middleware.JwtMiddleware{UcContract: route.Handler.UseCaseContract}
	isActiveStatusMiddleware := middleware.VerifyIsActiveStatus{UcContract: jwtMiddleware.UcContract}

	nurseRoutes := route.RouteGroup.Group("/nurse")
	nurseRoutes.Use(jwtMiddleware.New)
	nurseRoutes.Use(isActiveStatusMiddleware.IsActiveStatusMiddleware)
	nurseRoutes.Get("", handler.Browse)
	nurseRoutes.Get("/all", handler.BrowseAll)
	nurseRoutes.Get("/:id", handler.ReadByPk)
	nurseRoutes.Put("/:id", handler.Edit)
	nurseRoutes.Post("", handler.Add)
	nurseRoutes.Delete("/:id", handler.Delete)
}
