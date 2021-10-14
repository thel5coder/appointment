package routers

import (
	"github.com/gofiber/fiber/v2"
	"profira-backend/server/handlers"
)

type AuthenticationRoutes struct {
	RouteGroup fiber.Router
	Handler    handlers.Handler
}

func (route AuthenticationRoutes) RegisterRoute() {
	handler := handlers.AuthenticationHandler{Handler: route.Handler}

	authenticationRoutes := route.RouteGroup.Group("/auth")
	authenticationRoutes.Post("/login", handler.Login)
	authenticationRoutes.Post("/registration", handler.Registration)
	authenticationRoutes.Post("/forgot-password",handler.ForgotPassword)
	authenticationRoutes.Post("/customer-activation",handler.ActivationCustomer)
}
