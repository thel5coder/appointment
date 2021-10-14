package main

import (
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	idTranslations "github.com/go-playground/validator/v10/translations/id"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/rs/xid"
	"log"
	"time"

	"os"
	"profira-backend/server/bootstrap"
	"profira-backend/usecase"
)

var (
	validatorDriver *validator.Validate
	uni             *ut.UniversalTranslator
	translator      ut.Translator
	logFormat       = `{"host":"${host}","pid":"${pid}","time":"${time}","request-id":"${locals:requestid}","status":"${status}","method":"${method}","latency":"${latency}","path":"${path}",` +
		`"user-agent":"${ua}","in":"${bytesReceived}","out":"${bytesSent}"}`
)

func main() {
	config, err := bootstrap.LoadConfig()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer config.DB.Close()

	//init validator
	validatorInit()

	app := fiber.New()

	ucContract := usecase.UcContract{
		App:          app,
		ReqID:        xid.New().String(),
		DB:           config.DB,
		RedisClient:  config.RedisClient,
		Jwe:          config.JweCredential,
		Validate:     validatorDriver,
		Translator:   translator,
		JwtConfig:    config.JwtConfig,
		JwtCred:      config.JwtCredential,
		AWSS3:        config.AwsCredential,
		Pusher:       config.PusherCredential,
		GoMailConfig: config.MailingConfig,
		Fcm:          config.FcmConnection,
		AmqpConn:     config.AmqpConn,
		AmqpChannel:  config.AmqpChannel,
		TwilioClient: config.TwilioClient,
		Mandrill:     config.MandrillCred,
		MinioClient:  config.MinioClient,
	}

	boot := bootstrap.Bootstrap{
		App:             app,
		Db:              config.DB,
		UseCaseContract: ucContract,
		Jwe:             config.JweCredential,
		Translator:      translator,
		Validator:       validatorDriver,
		JwtConfig:       config.JwtConfig,
		JwtCred:         config.JwtCredential,
	}

	boot.App.Use(recover.New())
	boot.App.Use(requestid.New())
	boot.App.Use(cors.New())
	boot.App.Use(logger.New(logger.Config{
		Format:     logFormat + "\n",
		TimeFormat: time.RFC3339,
		TimeZone:   "Asia/Jakarta",
	}))

	boot.RegisterRouters()
	log.Fatal(boot.App.Listen(os.Getenv("APP_HOST_SERVER")))
}

func validatorInit() {
	en := en.New()
	id := id.New()
	uni = ut.New(en, id)

	transEN, _ := uni.GetTranslator("en")
	transID, _ := uni.GetTranslator("id")

	validatorDriver = validator.New()

	enTranslations.RegisterDefaultTranslations(validatorDriver, transEN)
	idTranslations.RegisterDefaultTranslations(validatorDriver, transID)

	switch os.Getenv("APP_LOCALE") {
	case "en":
		translator = transEN
	case "id":
		translator = transID
	}
}
