package server

import (
	"context"
	"net/http"
	"time"
	"websocket/database"
	"websocket/providers"
	"websocket/providers/chatProvider"
	"websocket/providers/dbHelperProvider"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Server struct {
	DB                   *sqlx.DB
	HttpServer           *http.Server
	RealtimeChatProvider providers.RealtimeChatHubProvider
	DBHelper             providers.DBHelperProvider
}

const (
	readTimeout       = 5 * time.Minute
	readHeaderTimeout = 30 * time.Second
	writeTimeout      = 5 * time.Minute
)

func ServerInit() *Server {

	psqldb := database.ConnectAndMigrate(
		"localhost",
		5432,
		"mydb",
		"shiv",
		"123",
		database.SSLModeDisable)

	logrus.Print("migration successful!!")

	realtimeChatProvider := chatProvider.NewRealtimeChatProvider()

	dbHelper := dbHelperProvider.NewDBHelperProvider(psqldb)

	return &Server{
		DB:                   psqldb,
		RealtimeChatProvider: realtimeChatProvider,
		DBHelper:             dbHelper,
	}
}

func (srv *Server) Start(port string) error {
	srv.HttpServer = &http.Server{
		Addr:              port,
		Handler:           srv.InjectRoutes(),
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
	}
	return srv.HttpServer.ListenAndServe()
}

func (srv *Server) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return srv.HttpServer.Shutdown(ctx)
}
