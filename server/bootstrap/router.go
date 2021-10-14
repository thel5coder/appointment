package bootstrap

import (
	"github.com/gofiber/fiber/v2"
	"profira-backend/server/bootstrap/routers"
	"profira-backend/server/handlers"
)

func (boot Bootstrap) RegisterRouters() {
	handlerType := handlers.Handler{
		App:             boot.App,
		UseCaseContract: &boot.UseCaseContract,
		Jwe:             boot.Jwe,
		Db:              boot.Db,
		Validate:        boot.Validator,
		Translator:      boot.Translator,
		JwtConfig:       boot.JwtConfig,
	}

	boot.App.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(200).JSON("Work")
	})

	apiV1 := boot.App.Group("/api/v1")

	//user routes
	userRoutes := routers.UserRoutes{
		RouteGroup: apiV1,
		Handler:    handlerType,
	}
	userRoutes.RegisterRoute()

	//otp routes
	otpRoutes := routers.OTPRoutes{
		RouteGroup: apiV1,
		Handler:    handlerType,
	}
	otpRoutes.RegisterRoute()

	//file routes
	fileRoutes := routers.FileRoutes{
		RouteGroup: apiV1,
		Handler:    handlerType,
	}
	fileRoutes.RegisterRoute()

	//authentication routes
	authenticationRoutes := routers.AuthenticationRoutes{
		RouteGroup: apiV1,
		Handler:    handlerType,
	}
	authenticationRoutes.RegisterRoute()

	//enum routes
	enumRoutes := routers.EnumRoutes{
		RouteGroup: apiV1,
		Handler:    handlerType,
	}
	enumRoutes.RegisterRoute()

	//role routes
	roleRoutes := routers.RoleRoutes{
		RouteGroup: apiV1,
		Handler:    handlerType,
	}
	roleRoutes.RegisterRoute()

	//admin routes
	adminRoutes := routers.AdminRoutes{
		RouteGroup: apiV1,
		Handler:    handlerType,
	}
	adminRoutes.RegisterRoute()

	//nurse routes
	nurseRoutes := routers.NurseRoutes{
		RouteGroup: apiV1,
		Handler:    handlerType,
	}
	nurseRoutes.RegisterRoute()

	//doctor routes
	doctorRoutes := routers.DoctorRoutes{
		RouteGroup: apiV1,
		Handler:    handlerType,
	}
	doctorRoutes.RegisterRoute()

	//clinic routes
	clinicRoutes := routers.ClinicRoutes{
		RouteGroup: apiV1,
		Handler:    handlerType,
	}
	clinicRoutes.RegisterRoute()

	//customer routes
	customerRoutes := routers.CustomerRoutes{
		RouteGroup: apiV1,
		Handler:    handlerType,
	}
	customerRoutes.RegisterRoute()

	//procedure routes
	procedureRoutes := routers.ProcedureRoutes{
		RouteGroup: apiV1,
		Handler:    handlerType,
	}
	procedureRoutes.RegisterRoute()

	//treatment routes
	treatmentRoutes := routers.TreatmentRoutes{
		RouteGroup: apiV1,
		Handler:    handlerType,
	}
	treatmentRoutes.RegisterRoute()

	//bed routes
	bedRoutes := routers.BedRoutes{
		RouteGroup: apiV1,
		Handler:    handlerType,
	}
	bedRoutes.RegisterRoute()

	////doctor week day schedule routes
	//doctorWeekDayScheduleRoutes := routers.DoctorWeekDayScheduleRoutes{
	//	RouteGroup: apiV1,
	//	Handler:    handlerType,
	//}
	//doctorWeekDayScheduleRoutes.RegisterRoute()

	//category routes
	categoryRoutes := routers.CategoryRoutes{
		RouteGroup: apiV1,
		Handler:    handlerType,
	}
	categoryRoutes.RegisterRoute()

	//promotion routes
	promotionRoutes := routers.PromotionRoutes{
		RouteGroup: apiV1,
		Handler:    handlerType,
	}
	promotionRoutes.RegisterRoute()

	//schedule doctor availability routes
	doctorAvailabilityRoutes := routers.DoctorAvailabilityRoutes{
		RouteGroup: apiV1,
		Handler:    handlerType,
	}
	doctorAvailabilityRoutes.RegisterRouter()
}
