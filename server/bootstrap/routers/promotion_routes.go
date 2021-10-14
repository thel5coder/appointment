package routers

import (
	"github.com/gofiber/fiber/v2"
	"profira-backend/server/handlers"
	"profira-backend/server/middleware"
)

type PromotionRoutes struct {
	RouteGroup fiber.Router
	Handler    handlers.Handler
}

func (route PromotionRoutes) RegisterRoute() {
	handler := handlers.PromotionHandler{Handler: route.Handler}
	jwtMiddleware := middleware.JwtMiddleware{UcContract: route.Handler.UseCaseContract}
	isActiveStatusMiddleware := middleware.VerifyIsActiveStatus{UcContract: jwtMiddleware.UcContract}

	promotionRoute := route.RouteGroup.Group("/promotion")
	promotionRoute.Use(jwtMiddleware.New)
	promotionRoute.Use(isActiveStatusMiddleware.IsActiveStatusMiddleware)
	promotionRoute.Get("", handler.Browse)
	promotionRoute.Get("/active", handler.BrowseActivePromotion)
	promotionRoute.Get("/:id", handler.ReadBy)
	promotionRoute.Put("/:id", handler.Edit)
	promotionRoute.Post("", handler.Add)
	promotionRoute.Delete("/:id", handler.Delete)
}
