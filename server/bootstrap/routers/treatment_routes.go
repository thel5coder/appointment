package routers

import (
	"github.com/gofiber/fiber/v2"
	"profira-backend/server/handlers"
	"profira-backend/server/middleware"
)

type TreatmentRoutes struct {
	RouteGroup fiber.Router
	Handler    handlers.Handler
}

func (route TreatmentRoutes) RegisterRoute() {
	handler := handlers.TreatmentHandler{Handler: route.Handler}
	jwtMiddleware := middleware.JwtMiddleware{UcContract: route.Handler.UseCaseContract}
	isActiveStatusMiddleware := middleware.VerifyIsActiveStatus{UcContract: jwtMiddleware.UcContract}

	treatmentRoutes := route.RouteGroup.Group("/treatment")
	treatmentRoutes.Use(jwtMiddleware.New)
	treatmentRoutes.Use(isActiveStatusMiddleware.IsActiveStatusMiddleware)
	treatmentRoutes.Get("", handler.Browse)
	treatmentRoutes.Get("/all", handler.BrowseAll)
	treatmentRoutes.Get("/category/:id", handler.BrowseByCategory)
	treatmentRoutes.Get("/:id", handler.ReadByPk)
	treatmentRoutes.Put("/:id", handler.Edit)
	treatmentRoutes.Post("", handler.Add)
	treatmentRoutes.Delete("/:id", handler.Delete)
}
