package routers

import (
	"github.com/gofiber/fiber/v2"
	"profira-backend/server/handlers"
	"profira-backend/server/middleware"
)

type CategoryRoutes struct {
	RouteGroup fiber.Router
	Handler    handlers.Handler
}

func (route CategoryRoutes) RegisterRoute() {
	handler := handlers.CategoryHandler{Handler: route.Handler}
	jwtMiddleware := middleware.JwtMiddleware{UcContract: route.Handler.UseCaseContract}
	//isActiveStatusMiddleware := middleware.VerifyIsActiveStatus{UcContract: jwtMiddleware.UcContract}

	categoryRoutes := route.RouteGroup.Group("/category")
	categoryRoutes.Use(jwtMiddleware.New)
	//categoryRoutes.Use(isActiveStatusMiddleware.IsActiveStatusMiddleware)
	categoryRoutes.Get("/tree-view", handler.BrowseTreeView)
	categoryRoutes.Get("/short-cut", handler.BrowseShortcutCategory)
	categoryRoutes.Get("/parent", handler.BrowseByParent)
	categoryRoutes.Get("/:id", handler.ReadByID)
	categoryRoutes.Put("/:id", handler.Edit)
	categoryRoutes.Post("", handler.Add)
	categoryRoutes.Delete("/:id", handler.Delete)
}
