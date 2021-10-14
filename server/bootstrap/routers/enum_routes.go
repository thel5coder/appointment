package routers

import (
	"github.com/gofiber/fiber/v2"
	"profira-backend/server/handlers"
	"profira-backend/server/middleware"
)

type EnumRoutes struct {
	RouteGroup fiber.Router
	Handler handlers.Handler
}

func (route EnumRoutes) RegisterRoute(){
	handler := handlers.EnumHandler{Handler:route.Handler}
	jwtMiddleware := middleware.JwtMiddleware{UcContract: route.Handler.UseCaseContract}

	enumRoutes := route.RouteGroup.Group("/enums")
	enumRoutes.Use(jwtMiddleware.New)
	enumRoutes.Get("/sex",handler.GetSexEnums)
}
