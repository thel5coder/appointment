package routers

import (
	"github.com/gofiber/fiber/v2"
	"profira-backend/server/handlers"
	"profira-backend/server/middleware"
)

type UserRoutes struct {
	RouteGroup fiber.Router
	Handler handlers.Handler
}

func(route UserRoutes) RegisterRoute(){
	handler := handlers.UserHandler{Handler: route.Handler}
	jwtMiddleware := middleware.JwtMiddleware{UcContract: route.Handler.UseCaseContract}
	isActiveStatusMiddleware := middleware.VerifyIsActiveStatus{UcContract: jwtMiddleware.UcContract}

	userRoutes := route.RouteGroup.Group("/user")
	userRoutes.Use(jwtMiddleware.New)
	userRoutes.Use(isActiveStatusMiddleware.IsActiveStatusMiddleware)
	userRoutes.Get("/profile",handler.ReadProfile)
}
