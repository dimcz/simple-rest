package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"simple-rest/data"
	"simple-rest/pkg/util"
	"simple-rest/routers"
	"time"
)

func main() {
	filename := flag.String("c", "app.yaml", "use an alternative configuration file")
	flag.Parse()

	logger := util.GetLogger()
	logger.Info("reading settings")
	config := util.GetConfig(*filename, logger)

	logger.Info("open connection to DB")
	data.NewConnection(config, logger)
	util.Setup(config)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d",
		config.BindIP, config.HTTPPort))
	if err != nil {
		logger.Fatal(err)
	}

	maxHeaderBytes := 1 << 20
	server := &http.Server{
		Handler:        routers.InitRouters(logger),
		ReadTimeout:    time.Second * 60,
		WriteTimeout:   time.Second * 60,
		MaxHeaderBytes: maxHeaderBytes,
	}
	logger.Infof("[info] start http server listening %s:%d", config.BindIP, config.HTTPPort)
	logger.Fatal(server.Serve(listener))
}
