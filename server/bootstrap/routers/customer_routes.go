package routers

import (
	"github.com/gofiber/fiber/v2"
	"profira-backend/server/handlers"
	"profira-backend/server/middleware"
)

type CustomerRoutes struct {
	RouteGroup fiber.Router
	Handler    handlers.Handler
}

func (route CustomerRoutes) RegisterRoute() {
	handler := handlers.CustomerHandler{Handler: route.Handler}
	jwtMiddleware := middleware.JwtMiddleware{UcContract: route.Handler.UseCaseContract}
	isActiveStatusMiddleware := middleware.VerifyIsActiveStatus{UcContract: jwtMiddleware.UcContract}

	customerRoutes := route.RouteGroup.Group("/customer")
	customerRoutes.Use(jwtMiddleware.New)
	customerRoutes.Use(isActiveStatusMiddleware.IsActiveStatusMiddleware)
	customerRoutes.Get("", handler.Browse)
	customerRoutes.Get("/profile", handler.MyProfile)
	customerRoutes.Post("/profile", handler.EditProfile)
	customerRoutes.Get("/:id", handler.ReadByPk)
	customerRoutes.Put("/:id", handler.Edit)
	customerRoutes.Post("", handler.Add)
	customerRoutes.Delete("/:id", handler.DeleteByPk)
}
