package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"simple-rest/logging"
	"simple-rest/model"
	"simple-rest/routers"
	"simple-rest/settings"
	"time"
)

func main() {

	logger := logging.GetLogger()

	filename := flag.String("c", "app.yaml", "use an alternative configuration file")
	migration := flag.Bool("s", true, "skip database migration")
	flag.Parse()

	logger.Info("reading settings")
	settings.Setup(*filename, logger)

	logger.Info("open connection to DB")
	model.Setup(*migration, logger)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d",
		settings.AppSettings.BindIP, settings.AppSettings.HTTPPort))
	if err != nil {
		logger.Fatal(err)
	}

	maxHeaderBytes := 1 << 20
	server := &http.Server{
		Handler:        routers.Setup(logger),
		ReadTimeout:    time.Second * 60,
		WriteTimeout:   time.Second * 60,
		MaxHeaderBytes: maxHeaderBytes,
	}
	logger.Infof("[info] start http server listening %s:%d", settings.AppSettings.BindIP, settings.AppSettings.HTTPPort)
	logger.Fatal(server.Serve(listener))
}
