package main

import (
	"log"
	"sync"

	"github.com/kiberlom/img_downloader/internal/background"
	"github.com/kiberlom/img_downloader/internal/config"
	"github.com/kiberlom/img_downloader/internal/db"
)

func main() {

	cnf := config.NewConfig()

	con, err := db.NewConnect(cnf)
	if err != nil {
		log.Fatal(err)
	}

	// запуск паука
	wg := &sync.WaitGroup{}
	wg.Add(1)
	background.SpiderUrl(con, wg)

	wg.Wait()

}
