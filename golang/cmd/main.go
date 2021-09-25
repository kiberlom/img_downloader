package main

import (
	"log"
	"sync"

	"github.com/kiberlom/img_downloader/internal/background"
	"github.com/kiberlom/img_downloader/internal/config"
	"github.com/kiberlom/img_downloader/internal/db"
	"github.com/kiberlom/img_downloader/internal/logger"
	"github.com/kiberlom/img_downloader/internal/shutdown"
)

func main() {

	sh := shutdown.NewShutdown()

	cnf := config.NewConfig()

	logf := logger.NewLogger()

	con, err := db.NewConnect(cnf)
	if err != nil {
		log.Fatal(err)
	}

	// запуск паука
	wgService := &sync.WaitGroup{}

	wgService.Add(1)

	go background.SpiderUrl(&background.ConfBackService{
		Shd: sh,
		WG:  wgService,
		Con: con,
		Log: logf,
	})

	go background.HostParse(&background.ConfBackService{
		Shd: sh,
		WG:  wgService,
		Con: con,
	})

	wgService.Wait()

	log.Println("Работа программы завершена корректно")

}
