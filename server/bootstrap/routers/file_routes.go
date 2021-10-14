package routers

import (
	"github.com/gofiber/fiber/v2"
	"profira-backend/server/handlers"
	"profira-backend/server/middleware"
)

type FileRoutes struct {
	RouteGroup fiber.Router
	Handler    handlers.Handler
}

func (route FileRoutes) RegisterRoute() {
	handler := handlers.FileHandler{Handler: route.Handler}
	jwtMiddleware := middleware.JwtMiddleware{UcContract: handler.UseCaseContract}

	fileRoutes := route.RouteGroup.Group("/file")
	fileRoutes.Use(jwtMiddleware.New)
	fileRoutes.Post("", handler.Add)
	fileRoutes.Get("/:id", handler.Read)
}
