package routers

import (
	"github.com/gofiber/fiber/v2"
	"profira-backend/server/handlers"
	"profira-backend/server/middleware"
)

type BedRoutes struct {
	RouteGroup fiber.Router
	Handler    handlers.Handler
}

func (route BedRoutes) RegisterRoute() {
	handler := handlers.BedHandler{Handler: route.Handler}
	jwtMiddleware := middleware.JwtMiddleware{UcContract: route.Handler.UseCaseContract}
	isActiveStatusMiddleware := middleware.VerifyIsActiveStatus{UcContract: jwtMiddleware.UcContract}

	bedRoutes := route.RouteGroup.Group("/bed")
	bedRoutes.Use(jwtMiddleware.New)
	bedRoutes.Use(isActiveStatusMiddleware.IsActiveStatusMiddleware)
	bedRoutes.Get("", handler.Browse)
	bedRoutes.Get("/all", handler.BrowseAll)
	bedRoutes.Get("/:id", handler.ReadByPk)
	bedRoutes.Put("/:id", handler.Edit)
	bedRoutes.Post("", handler.Add)
	bedRoutes.Delete("/:id", handler.DeleteByPk)
}
