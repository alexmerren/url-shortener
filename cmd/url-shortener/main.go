package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"url-shortner/internal/config"
	"url-shortner/internal/datastore"
	"url-shortner/internal/logger"
	"url-shortner/internal/rest"
)

const (
	ErrAppBuild     = 1
	ErrServerStop   = 2
	ErrLoggerSync   = 3
	ErrDatabaseStop = 4
)

type Application struct {
	server   *rest.RESTServer
	database datastore.UrlStorer
	logger   logger.Logger
}

func NewApplication() *Application {
	filesystem := config.NewFilesystem()
	config := config.NewConfiguration("config.yaml", filesystem)

	logLevel, err := config.GetString("logger.level")
	logEncoding, err := config.GetString("logger.encoding")
	logger, err := logger.NewZapLogger(logEncoding, logLevel)
	if err != nil {
		log.Print(err)
		return nil
	}

	ctx := context.Background()
	dbHost, err := config.GetString("database.host")
	dbPort, err := config.GetInt("database.port")
	dbUser, err := config.GetString("database.user")
	dbPassword, err := config.GetString("database.password")
	dbName, err := config.GetString("database.name")
	dbCapacity, err := config.GetInt("database.capacity")
	if err != nil {
		logger.WithError(err).Error("could not get database config")
	}

	database := datastore.NewPostgresUrlStore(
		ctx, logger, dbUser, dbPassword, dbHost, dbName, dbPort, dbCapacity,
	)

	restHost, err := config.GetString("rest.host")
	restPort, err := config.GetInt("rest.port")
	if err != nil {
		logger.WithError(err).Error("could not get rest config")
		return nil
	}

	router := rest.NewRouter(logger, database)

	restServer := rest.NewRESTServer(
		logger, database, restHost, restPort, router,
	)

	return &Application{
		server:   restServer,
		logger:   logger,
		database: database,
	}
}

func (a *Application) Run() int {
	terminationChannel := make(chan os.Signal, 1)
	signal.Notify(terminationChannel, syscall.SIGINT, syscall.SIGTERM)

	go a.server.Start()
	go a.database.Start()

	<-terminationChannel

	if err := a.server.Stop(); err != nil {
		return ErrServerStop
	}

	if err := a.logger.Sync(); err != nil {
		return ErrLoggerSync
	}

	if err := a.database.Stop(); err != nil {
		return ErrDatabaseStop
	}

	return 0
}

func main() {
	app := NewApplication()
	if app == nil {
		os.Exit(ErrAppBuild)
	}
	if code := app.Run(); code != 0 {
		os.Exit(code)
	}
}
