package server

import (
	"net/http"
	"strconv"

	"VivekPapnaiAtRS/template/providers/realtimesocketmanager"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	defaultReadWriteBufferSize = 1024
)

var realtimeUpgrader = websocket.Upgrader{
	ReadBufferSize:  defaultReadWriteBufferSize,
	WriteBufferSize: defaultReadWriteBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (srv *Server) realTimeWS(resp http.ResponseWriter, req *http.Request) {
	//uc := srv.getUserContext(req)
	conn, err := realtimeUpgrader.Upgrade(resp, req, nil)
	if err != nil {
		logrus.Error(err)
		return
	}

	userID, err := strconv.Atoi(req.URL.Query().Get("userId"))
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
	}

	client := realtimesocketmanager.NewRealtimeClient(srv.RealtimeHub.Get().(*realtimesocketmanager.RealtimeHub), conn, userID)

	client.Register()
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
	go client.ReadPump()
}
