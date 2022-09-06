package server

import (
	"context"
	"net/http"
	"os"
	"time"

	"VivekPapnaiAtRS/template/providers"
	"VivekPapnaiAtRS/template/providers/dbhelpprovider"
	"VivekPapnaiAtRS/template/providers/dbprovider"
	"VivekPapnaiAtRS/template/providers/realtimesocketmanager"
	"github.com/sirupsen/logrus"
)

const (
	defaultServerRequestTimeoutMinutes      = 2
	defaultServerReadHeaderTimeoutSeconds   = 30
	defaultServerRequestWriteTimeoutMinutes = 30
)

type Server struct {
	PSQL        providers.PSQLProvider
	DBHelper    providers.DBHelpProvider
	RealtimeHub providers.WebSocketHubProvider
	httpServer  *http.Server
}

func SrvInit() *Server {
	// PSQL connection
	db := dbprovider.NewPSQLProvider(os.Getenv("PSQL_STRING"), 10, 10)

	dbHelper := dbhelpprovider.NewDBHelper(db.DB())

	// realtimeHub is broadcaster for all other realtime messages
	realtimeHub := realtimesocketmanager.NewRealtimeHub(db.DB(), dbHelper)

	return &Server{
		PSQL:        db,
		DBHelper:    dbHelper,
		RealtimeHub: realtimeHub,
	}
}

func (srv *Server) Start() {
	addr := ":" + os.Getenv("PORT")

	httpSrv := &http.Server{
		Addr:              addr,
		Handler:           srv.InjectRoutes(),
		ReadTimeout:       defaultServerRequestTimeoutMinutes * time.Minute,
		ReadHeaderTimeout: defaultServerReadHeaderTimeoutSeconds * time.Second,
		WriteTimeout:      defaultServerRequestWriteTimeoutMinutes * time.Minute,
	}
	srv.httpServer = httpSrv

	go srv.RealtimeHub.Run()

	logrus.Info("Server running at PORT ", addr)
	if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.Fatalf("Start %v", err)
		return
	}
}

func (srv *Server) Stop() {
	logrus.Info("closing Postgres...")
	_ = srv.PSQL.DB().Close()

	logrus.Info("closing Websocket RealtimeHub...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logrus.Info("closing server...")
	_ = srv.httpServer.Shutdown(ctx)
	logrus.Info("Done")
}
