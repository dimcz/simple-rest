package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"net/http"
	"simple-rest/model"
	"simple-rest/routers"
	"simple-rest/settings"
	"time"
)

func main() {
	filename := flag.String("c", "app.yaml", "use an alternative configuration file")
	migration := flag.Bool("s", true, "skip database migration")
	flag.Parse()

	settings.Setup(*filename)
	model.Setup(*migration)

	gin.SetMode(settings.AppSettings.ServerMode)

	maxHeaderBytes := 1 << 20

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d",
		settings.AppSettings.BindIP, settings.AppSettings.HTTPPort))
	if err != nil {
		log.Fatalf("main.Listen error: %v", err)
	}

	server := &http.Server{
		Handler:        routers.Setup(),
		ReadTimeout:    time.Second * 60,
		WriteTimeout:   time.Second * 60,
		MaxHeaderBytes: maxHeaderBytes,
	}
	log.Printf("[info] start http server listening %s:%d", settings.AppSettings.BindIP, settings.AppSettings.HTTPPort)
	log.Fatal(server.Serve(listener))
}
