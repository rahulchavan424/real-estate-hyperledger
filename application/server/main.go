package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"application/blockchain"
	"application/pkg/cron"
	"application/routers"
)

func main() {
	timeLocal, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Printf("Time zone setting failed: %s", err)
	}
	time.Local = timeLocal

	blockchain.Init()
	go cron.Init()

	endPoint := fmt.Sprintf("0.0.0.0:%d", 8888)
	server := &http.Server{
		Addr:    endPoint,
		Handler: routers.InitRouter(),
	}
	log.Printf("[info] Start HTTP server listening on %s", endPoint)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Start HTTP server failed: %s", err)
	}
}
