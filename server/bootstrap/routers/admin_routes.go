package routers

import (
	"github.com/gofiber/fiber/v2"
	"profira-backend/server/handlers"
	"profira-backend/server/middleware"
)

type AdminRoutes struct {
	RouteGroup fiber.Router
	Handler    handlers.Handler
}

func (route AdminRoutes) RegisterRoute() {
	handler := handlers.AdminHandler{Handler: route.Handler}
	jwtMiddleware := middleware.JwtMiddleware{UcContract: route.Handler.UseCaseContract}
	isActiveStatusMiddleware := middleware.VerifyIsActiveStatus{UcContract: jwtMiddleware.UcContract}

	adminRoutes := route.RouteGroup.Group("/user-admin")
	adminRoutes.Use(jwtMiddleware.New)
	adminRoutes.Use(isActiveStatusMiddleware.IsActiveStatusMiddleware)
	adminRoutes.Get("", handler.Browse)
	adminRoutes.Get("/all", handler.BrowseAll)
	adminRoutes.Get("/:id", handler.ReadByPk)
	adminRoutes.Put("/:id", handler.Edit)
	adminRoutes.Post("", handler.Add)
	adminRoutes.Delete("/:id", handler.DeleteByPk)
}
