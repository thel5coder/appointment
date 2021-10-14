package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rs/xid"
	"github.com/streadway/amqp"
	"log"
	"os"
	amqpHelper "profira-backend/helpers/amqp"
	amqpConsumerHelper "profira-backend/helpers/amqpconsumer"
	"profira-backend/helpers/interfacepkg"
	"profira-backend/helpers/logruslogger"
	"profira-backend/server/bootstrap"
	"profira-backend/usecase"
	"runtime"
	"strconv"
)

var (
	uri          *string
	formURL      = flag.String("form_url", "http://localhost", "The URL that requests are sent to")
	logFile      = flag.String("log_file", "system.log", "The file where errors are logged")
	threads      = flag.Int("threads", 1, "The max amount of go routines that you would like the process to use")
	maxprocs     = flag.Int("max_procs", 1, "The max amount of processors that your application should use")
	paymentsKey  = flag.String("payments_key", "secret", "Access key")
	exchange     = flag.String("exchange", amqpHelper.MailExchange, "The exchange we will be binding to")
	exchangeType = flag.String("exchange_type", "direct", "Type of exchange we are binding to | topic | direct| etc..")
	queue        = flag.String("queue", amqpHelper.MailIncoming, "Name of the queue that you would like to connect to")
	routingKey   = flag.String("routing_key", amqpHelper.MailDeadLetter, "queue to route messages to")
	workerName   = flag.String("worker_name", "worker.name", "name to identify worker by")
	verbosity    = flag.Bool("verbos", false, "Set true if you would like to log EVERYTHING")

	// Hold consumer so our go routine can listen to
	// it's done error chan and trigger reconnects
	// if it's ever returned
	conn *amqpConsumerHelper.Consumer
)

func init() {
	flag.Parse()
	runtime.GOMAXPROCS(*maxprocs)
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("Error loading ..env file")
	}
	uri = flag.String("uri", os.Getenv("AMQP_URL"), "The rabbitmq endpoint")
}

func main() {
	file := false
	if file {
		f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		defer f.Close()
		if err != nil {
			panic(err)
		}
		log.SetOutput(f)
	}

	conn := amqpConsumerHelper.NewConsumer(*workerName, *uri, *exchange, *exchangeType, *queue)
	if err := conn.Connect(); err != nil {
		panic(err)
	}

	deliveries, err := conn.AnnounceQueue(*queue, *routingKey)
	if err != nil {
		panic(err)
	}

	config, err := bootstrap.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	ucContract := usecase.UcContract{
		ReqID:        xid.New().String(),
		DB:           config.DB,
		RedisClient:  config.RedisClient,
		Jwe:          config.JweCredential,
		JwtConfig:    config.JwtConfig,
		JwtCred:      config.JwtCredential,
		GoMailConfig: config.MailingConfig,
		Fcm:          config.FcmConnection,
		Mandrill:     config.MandrillCred,
	}

	conn.Handle(deliveries, handler, *threads, *queue, *routingKey, ucContract)
}

func handler(deliveries <-chan amqp.Delivery, uc *usecase.UcContract) {
	ctx := "sendingMailListener"
	for d := range deliveries {
		var formData map[string]interface{}
		var payload map[string]interface{}

		err := json.Unmarshal(d.Body, &formData)
		if err != nil {
			log.Printf("Error unmarshaling data: %s", err.Error())
		}

		payloadByte,err := json.Marshal(formData["payload"])
		err = json.Unmarshal(payloadByte,&payload)
		if err != nil {
			log.Printf("Error unmarshaling data: %s", err.Error())
		}

		logruslogger.Log(logruslogger.InfoLevel, interfacepkg.Marshall(formData), ctx, "begin", formData["qid"].(string))

		uc.ReqID = formData["qid"].(string)
		mailUc := usecase.MailUseCase{UcContract: uc}
		res, err := mailUc.SendMail(payload)
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "err", formData["qid"].(string))

			// Get fail counter from redis
			failCounter := amqpConsumerHelper.FailCounter{}
			err = uc.RedisClient.GetFromRedis("amqpFail"+formData["qid"].(string), &failCounter)
			if err != nil {
				failCounter = amqpConsumerHelper.FailCounter{
					Counter: 1,
				}
			}

			if failCounter.Counter > amqpConsumerHelper.MaxFailCounter {
				logruslogger.Log(logruslogger.WarnLevel, strconv.Itoa(failCounter.Counter), ctx, "rejected", formData["qid"].(string))
				d.Reject(false)
			} else {
				// Save the new counter to redis
				failCounter.Counter++
				err = uc.RedisClient.StoreToRedistWithExpired("amqpFail"+formData["qid"].(string), failCounter, "10m")

				logruslogger.Log(logruslogger.WarnLevel, strconv.Itoa(failCounter.Counter), ctx, "failed", formData["qid"].(string))
				d.Nack(false, true)
			}
		} else {
			logruslogger.Log(logruslogger.InfoLevel, interfacepkg.Marshall(res), ctx, "success", formData["qid"].(string))
			d.Ack(false)
		}
	}

	return
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
