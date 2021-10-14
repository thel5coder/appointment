package routers

import (
	"github.com/gofiber/fiber/v2"
	"profira-backend/server/handlers"
	"profira-backend/server/middleware"
)

type ProcedureRoutes struct {
	RouteGroup fiber.Router
	Handler    handlers.Handler
}

func (route ProcedureRoutes) RegisterRoute() {
	handler := handlers.ProcedureHandler{Handler: route.Handler}
	jwtMiddleware := middleware.JwtMiddleware{UcContract: route.Handler.UseCaseContract}
	isActiveStatusMiddleware := middleware.VerifyIsActiveStatus{UcContract: jwtMiddleware.UcContract}

	procedureRoutes := route.RouteGroup.Group("/procedure")
	procedureRoutes.Use(jwtMiddleware.New)
	procedureRoutes.Use(isActiveStatusMiddleware.IsActiveStatusMiddleware)
	procedureRoutes.Get("", handler.Browse)
	procedureRoutes.Get("/all", handler.BrowseAll)
	procedureRoutes.Get("/:id", handler.ReadByPk)
	procedureRoutes.Put("/:id", handler.Edit)
	procedureRoutes.Post("", handler.Add)
	procedureRoutes.Delete("/:id", handler.DeleteByPk)
}
