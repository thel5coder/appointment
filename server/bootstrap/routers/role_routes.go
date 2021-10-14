package routers

import (
	"github.com/gofiber/fiber/v2"
	"profira-backend/server/handlers"
	"profira-backend/server/middleware"
)

type RoleRoutes struct {
	RouteGroup fiber.Router
	Handler    handlers.Handler
}

func (route RoleRoutes) RegisterRoute() {
	handler := handlers.RoleHandler{Handler: route.Handler}
	jwtMiddleware := middleware.JwtMiddleware{UcContract: route.Handler.UseCaseContract}
	isActiveStatusMiddleware := middleware.VerifyIsActiveStatus{UcContract: jwtMiddleware.UcContract}

	roleRoute := route.RouteGroup.Group("/role")
	roleRoute.Use(jwtMiddleware.New)
	roleRoute.Use(isActiveStatusMiddleware.IsActiveStatusMiddleware)
	roleRoute.Get("", handler.Browse)
	roleRoute.Get("/:id", handler.Read)
	roleRoute.Put("/:id", handler.Edit)
	roleRoute.Post("", handler.Add)
	roleRoute.Delete("/:id", handler.Delete)
}
