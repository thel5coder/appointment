package bootstrap

import (
	"database/sql"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/go-redis/redis/v7"
	jwtFiber "github.com/gofiber/jwt/v2"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/streadway/amqp"
	"log"
	"os"
	"profira-backend/db"
	amqpHelper "profira-backend/helpers/amqp"
	awsHelper "profira-backend/helpers/aws"
	fcmHelper "profira-backend/helpers/fcm"
	"profira-backend/helpers/jwe"
	jwtHelper "profira-backend/helpers/jwt"
	mailingHelper "profira-backend/helpers/mailing"
	"profira-backend/helpers/mandrill"
	minioHelper "profira-backend/helpers/minio"
	pusherHelper "profira-backend/helpers/pusher"
	redisHelper "profira-backend/helpers/redis"
	"profira-backend/helpers/str"
	twilioHelper "profira-backend/helpers/twilio"
	"profira-backend/usecase"
)

type Config struct {
	UcContract       usecase.UcContract
	JweCredential    jwe.Credential
	RedisClient      redisHelper.RedisClient
	JwtConfig        jwtFiber.Config
	JwtCredential    jwtHelper.JwtCredential
	DB               *sql.DB
	AwsCredential    awsHelper.AWSS3
	PusherCredential pusherHelper.Credential
	MailingConfig    mailingHelper.GoMailConfig
	FcmConnection    fcmHelper.Connection
	AmqpConn         *amqp.Connection
	AmqpChannel      *amqp.Channel
	TwilioClient     *twilioHelper.Client
	MandrillCred     mandrill.Credential
	MinioClient      *minio.Client
}

func LoadConfig() (res Config, err error) {
	err = godotenv.Load("../.env")
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("Error loading ..env file")
	}

	//jwe
	res.JweCredential = jwe.Credential{
		KeyLocation: os.Getenv("PRIVATE_KEY"),
		Passphrase:  os.Getenv("PASSPHRASE"),
	}

	//setup redis
	redisOption := &redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	}
	res.RedisClient = redisHelper.RedisClient{Client: redis.NewClient(redisOption)}

	//jwtconfig
	res.JwtConfig = jwtFiber.Config{
		SigningKey: []byte(os.Getenv("SECRET")),
		Claims:     &jwtHelper.CustomClaims{},
	}

	//jwt credential
	res.JwtCredential = jwtHelper.JwtCredential{
		TokenSecret:         os.Getenv("SECRET"),
		ExpiredToken:        str.StringToInt(os.Getenv("TOKEN_EXP_TIME")),
		RefreshTokenSecret:  os.Getenv("SECRET_REFRESH_TOKEN"),
		ExpiredRefreshToken: str.StringToInt(os.Getenv("REFRESH_TOKEN_EXP_TIME")),
	}

	//twillio
	twilioClient := twilioHelper.NewTwilioClient(os.Getenv("TWILLIO_SID"), os.Getenv("TWILLIO_TOKEN"), os.Getenv("TWILLIO_SEND_FROM"))
	res.TwilioClient = twilioClient

	//mandril
	res.MandrillCred = mandrill.Credential{
		Key:      os.Getenv("MANDRILL_KEY"),
		FromMail: os.Getenv("MANDRILL_FROM_MAIL"),
		FromName: os.Getenv("MANDRILL_FROM_NAME"),
	}

	//setup db connection
	dbInfo := db.Connection{
		Host:                    os.Getenv("DB_HOST"),
		DbName:                  os.Getenv("DB_NAME"),
		User:                    os.Getenv("DB_USER"),
		Password:                os.Getenv("DB_PASSWORD"),
		Port:                    os.Getenv("DB_PORT"),
		SslMode:                 os.Getenv("DB_SSL_MODE"),
		DBMaxConnection:         str.StringToInt(os.Getenv("DB_MAX_CONNECTION")),
		DBMAxIdleConnection:     str.StringToInt(os.Getenv("DB_MAX_IDLE_CONNECTION")),
		DBMaxLifeTimeConnection: str.StringToInt(os.Getenv("DB_MAX_LIFETIME_CONNECTION")),
	}
	res.DB, err = dbInfo.DbConnect()
	if err != nil {
		panic(err)
	}

	pong, err := res.RedisClient.Client.Ping().Result()
	fmt.Println("Redis ping status: "+pong, err)

	//aws setup
	awsAccessKey := os.Getenv("S3_ACCESS_KEY")
	awsSecretKey := os.Getenv("S3_SECRET_KEY")
	awsBucket := os.Getenv("S3_BUCKET")
	awsDirectory := os.Getenv("S3_DIRECTORY")
	s3EndPoint := os.Getenv("S3_BASE_URL")
	awsConfig := aws.Config{
		Endpoint:    &s3EndPoint,
		Region:      aws.String(os.Getenv("S3_LOCATION")),
		Credentials: credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, ""),
	}
	res.AwsCredential = awsHelper.AWSS3{
		AWSConfig: awsConfig,
		Bucket:    awsBucket,
		Directory: awsDirectory,
		AccessKey: awsAccessKey,
		SecretKey: awsSecretKey,
	}

	//pusher
	res.PusherCredential = pusherHelper.Credential{
		AppID:   os.Getenv("PUSHER_APP_ID"),
		Key:     os.Getenv("PUSHER_KEY"),
		Secret:  os.Getenv("PUSHER_SECRET"),
		Cluster: os.Getenv("PUSHER_CLUSTER"),
	}

	//gomail config
	res.MailingConfig = mailingHelper.GoMailConfig{
		SMTPHost: os.Getenv("SMTP_HOST"),
		SMTPPort: str.StringToInt(os.Getenv("SMTP_PORT")),
		Sender:   os.Getenv("MAIL_SENDER"),
		Password: os.Getenv("PASSWORD"),
	}

	// FCM connection
	res.FcmConnection = fcmHelper.Connection{
		APIKey: os.Getenv("FMC_KEY"),
	}

	// Mqueue connection
	amqpInfo := amqpHelper.Connection{
		URL: os.Getenv("AMQP_URL"),
	}
	amqpConn, amqpChannel, err := amqpInfo.Connect()
	if err != nil {
		panic(err)
	}
	res.AmqpConn = amqpConn
	res.AmqpChannel = amqpChannel

	minioConn := minioHelper.Connection{
		AccessKey: os.Getenv("MINIO_ACCESS_KEY"),
		SecretKey: os.Getenv("MINIO_SECRET_KEY"),
		UseSSL:    os.Getenv("MINIO_USE_SSL"),
		EndPoint:  os.Getenv("MINIO_ENDPOINT"),
	}
	res.MinioClient, err = minioConn.InitClient()

	return res, err
}
