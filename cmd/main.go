package main

import (
	"VivekPapnaiAtRS/template/server"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("connected")
	srv := server.SrvInit()
	go srv.Start()
	<-done
	logrus.Info("Graceful shutdown")
	srv.Stop()
}
