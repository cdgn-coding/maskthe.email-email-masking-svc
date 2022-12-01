package http

import (
	"email-masking-svc/src/application/consumers"
	"email-masking-svc/src/application/services"
	"email-masking-svc/src/infrastructure/configuration"
	"email-masking-svc/src/infrastructure/http/controllers"
	"email-masking-svc/src/infrastructure/postgresql"
	"email-masking-svc/src/infrastructure/postgresql/repositories"
	"email-masking-svc/src/infrastructure/rabbitmq"
	"email-masking-svc/src/infrastructure/rabbitmq/queues"
	"net/http"

	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Server struct {
	config               configuration.Config
	logger               configuration.Logger
	createMaskController *controllers.CreateEmailMaskController
	rabbitmqChannel      *amqp.Channel
	emailConsumer        *consumers.EmailConsumer
}

func NewServer() *Server {
	config := configuration.LoadConfig()
	loggerLevel := config.GetString("logger.level")
	logger := configuration.NewLogger(loggerLevel)

	postgresConnection := postgresql.NewConnection(config.GetString("postgres.dsn"))
	maskRepository := repositories.NewPostgresMaskRepository(postgresConnection)
	createMask := services.NewCreateMask(maskRepository)
	createMaskController := controllers.NewCreateEmailMaskController(createMask, logger)

	channel := rabbitmq.CreateConnection(config.GetString("rabbitmq.url"))
	emailsToSendTopic := config.GetString("rabbitmq.queues.emailsToSend")
	sendEmailPublisher := queues.NewPublisher(channel, emailsToSendTopic)
	redirectEmail := services.NewRedirectEmail(maskRepository, sendEmailPublisher, logger)
	emailConsumer := consumers.NewEmailConsumer(logger, redirectEmail)

	return &Server{
		config:               config,
		logger:               logger,
		createMaskController: createMaskController,
		rabbitmqChannel:      channel,
		emailConsumer:        emailConsumer,
	}
}

func (s Server) Run() {
	s.logger.Info("Starting up server")
	s.bindConsumers()
	s.bindRoutes()
}

func (s Server) bindConsumers() {
	receivedEmailsTopic := s.config.GetString("rabbitmq.queues.emailsReceived")
	queues.ConsumeQueue(s.rabbitmqChannel, receivedEmailsTopic, s.emailConsumer)
}

func (s Server) bindRoutes() {
	router := mux.NewRouter()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("healthy"))
	})

	router.HandleFunc("/masks", s.createMaskController.ServeHTTP).Methods("POST")

	s.logger.Fatal(http.ListenAndServe(s.config.GetString("http.port"), router))
}
