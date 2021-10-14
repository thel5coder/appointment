package usecase

import (
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
	jwtFiber "github.com/gofiber/jwt/v2"
	"github.com/minio/minio-go/v7"
	"math/rand"
	"os"
	queue "profira-backend/helpers/amqp"
	"profira-backend/helpers/aws"
	"profira-backend/helpers/fcm"
	"profira-backend/helpers/jwe"
	"profira-backend/helpers/jwt"
	"profira-backend/helpers/logruslogger"
	"profira-backend/helpers/mailing"
	"profira-backend/helpers/mandrill"
	"profira-backend/helpers/messages"
	"profira-backend/helpers/pusher"
	twilioHelper "profira-backend/helpers/twilio"
	"profira-backend/usecase/viewmodel"
	"strings"
	"time"

	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	"github.com/streadway/amqp"
	"profira-backend/helpers/redis"
)

const (
	defaultLimit   = 10
	maxLimit       = 50
	defaultOrderBy = "id"
	defaultSort    = "asc"
	// PasswordLength ...
	PasswordLength  = 6
	defaultLastPage = 0
	// OtpLifeTime ...
	OtpLifeTime = "3m"
	// MaxOtpSubmitRetry ...
	MaxOtpSubmitRetry     = 3
	defaultCustomerRoleID = "df784016-1bb6-4069-a0ad-12936a09d87d"
	defaultSex            = "other"
)

// GlobalSmsCounter ...
var GlobalSmsCounter int

// AmqpConnection ...
var AmqpConnection *amqp.Connection

// AmqpChannel ...
var AmqpChannel *amqp.Channel

//X-Request-ID
var xRequestID interface{}

// UcContract ...
type UcContract struct {
	ReqID        string
	App          *fiber.App
	DB           *sql.DB
	TX           *sql.Tx
	AmqpConn     *amqp.Connection
	AmqpChannel  *amqp.Channel
	RedisClient  redis.RedisClient
	Jwe          jwe.Credential
	Validate     *validator.Validate
	Translator   ut.Translator
	JwtConfig    jwtFiber.Config
	JwtCred      jwt.JwtCredential
	AWSS3        aws.AWSS3
	Pusher       pusher.Credential
	GoMailConfig mailing.GoMailConfig
	UserID       string
	ClinicID     string
	Fcm          fcm.Connection
	TwilioClient *twilioHelper.Client
	Mandrill     mandrill.Credential
	MinioClient  *minio.Client
}

func (uc UcContract) setPaginationParameter(page, limit int, order, sort string) (int, int, int, string, string) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > maxLimit {
		limit = defaultLimit
	}
	if order == "" {
		order = defaultOrderBy
	}
	if sort == "" {
		sort = defaultSort
	}
	offset := (page - 1) * limit

	return offset, limit, page, order, sort
}

func (uc UcContract) setPaginationResponse(page, limit, total int) (paginationResponse viewmodel.PaginationVm) {
	var lastPage int

	if total > 0 {
		lastPage = total / limit

		if total%limit != 0 {
			lastPage = lastPage + 1
		}
	} else {
		lastPage = defaultLastPage
	}

	paginationResponse = viewmodel.PaginationVm{
		CurrentPage: page,
		LastPage:    lastPage,
		Total:       total,
		PerPage:     limit,
	}

	return paginationResponse
}

// GetRandomString ...
func (uc UcContract) GetRandomString(length int) string {
	if length == 0 {
		length = PasswordLength
	}

	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÅÄÖ" +
		"abcdefghijklmnopqrstuvwxyzåäö" +
		"0123456789")

	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	password := b.String()

	return password
}

// LimitRetryByKey ...
func (uc UcContract) LimitRetryByKey(key string, limit float64) (err error) {
	var count float64
	res := map[string]interface{}{}

	err = uc.RedisClient.GetFromRedis(key, &res)
	if err != nil {
		err = nil
		res = map[string]interface{}{
			"counter": count,
		}
	}
	count = res["counter"].(float64) + 1
	if count > limit {
		uc.RedisClient.RemoveFromRedis(key)

		return errors.New(messages.MaxRetryKey)
	}

	res["counter"] = count
	uc.RedisClient.StoreToRedistWithExpired(key, res, "24h")

	return err
}

// PushToQueue ...
func (uc UcContract) PushToQueue(queueBody map[string]interface{}, queueType, deadLetterType string) (err error) {
	mqueue := queue.NewQueue(AmqpConnection, AmqpChannel)

	_, _, err = mqueue.PushQueueReconnect(os.Getenv("AMQP_URL"), queueBody, queueType, deadLetterType)
	if err != nil {
		return err
	}

	return err
}

// InitDBTransaction ...
func (uc UcContract) InitDBTransaction() (err error) {
	uc.TX, err = uc.DB.Begin()
	if err != nil {
		uc.TX.Rollback()

		return err
	}

	return nil
}

func (uc UcContract) ErrorHandler(ctx, scope string, err error) error {
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, scope, uc.ReqID)
		return err
	}

	return nil
}
