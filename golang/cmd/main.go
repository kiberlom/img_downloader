package main

import (
	"fmt"
	"log"
	"time"

	"github.com/kiberlom/img_downloader/internal/background"
	"github.com/kiberlom/img_downloader/internal/config"
	"github.com/kiberlom/img_downloader/internal/db"
	"github.com/kiberlom/img_downloader/internal/geturl"
	"github.com/kiberlom/img_downloader/internal/parslink"
)

func rrr() {
	// u := "https://img5.goodfon.ru/original/2048x1365/4/ea/park-doroga-utro.jpg"

	// content, err := geturl.GetImage(u)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if err := file.Save(*content); err != nil {
	// 	log.Fatal(err)
	// }

	u := []string{"https://img5.goodfon.ru",
		"https://play.google.com/store/apps/details?id=com.goodfon.goodfon&rdid=com.goodfon.goodfon",
		"https://coderoad.ru/1821811/%D0%9A%D0%B0%D0%BA-%D1%87%D0%B8%D1%82%D0%B0%D1%82%D1%8C-%D0%BF%D0%B8%D1%81%D0%B0%D1%82%D1%8C-%D0%B8%D0%B7-%D0%B2-%D1%84%D0%B0%D0%B9%D0%BB-%D1%81-%D0%BF%D0%BE%D0%BC%D0%BE%D1%89%D1%8C%D1%8E-Go",
	}

	for _, v := range u {
		html, err := geturl.GetHtml(v)
		if err != nil {
			log.Fatal(err)
		}

		//fmt.Println(html)

		links, err := parslink.GetAllUrl(html)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(links)
	}

	for {
		time.Sleep(2 * time.Second)
		fmt.Println("Я тут...")
	}
}

func main() {

	cnf := config.NewConfig()

	con, err := db.NewConnect(cnf)
	if err != nil {
		log.Fatal(err)
	}

	// запуск паука
	background.SpiderUrl(con)

	r, err := con.FreeUrl()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%+v", r)

	// for {
	// 	time.Sleep(3 * time.Second)
	// 	test, err := con.Ping()
	// 	if err != nil {
	// 		log.Println(err)
	// 		continue
	// 	}

	// 	if test {
	// 		log.Println("Есть Ping")
	// 		continue
	// 	}

	// 	log.Println("Нет Pinga")

	// }

}
