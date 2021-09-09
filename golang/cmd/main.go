package main

import (
	"fmt"
	"log"

	"github.com/kiberlom/img_downloader/internal/geturl"
	"github.com/kiberlom/img_downloader/internal/parslink"
)

func main() {
	// u := "https://img5.goodfon.ru/original/2048x1365/4/ea/park-doroga-utro.jpg"

	// content, err := geturl.GetImage(u)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if err := file.Save(*content); err != nil {
	// 	log.Fatal(err)
	// }

	u := "https://img5.goodfon.ru"
	//u := "https://play.google.com/store/apps/details?id=com.goodfon.goodfon&rdid=com.goodfon.goodfon"
	// u := "https://coderoad.ru/1821811/%D0%9A%D0%B0%D0%BA-%D1%87%D0%B8%D1%82%D0%B0%D1%82%D1%8C-%D0%BF%D0%B8%D1%81%D0%B0%D1%82%D1%8C-%D0%B8%D0%B7-%D0%B2-%D1%84%D0%B0%D0%B9%D0%BB-%D1%81-%D0%BF%D0%BE%D0%BC%D0%BE%D1%89%D1%8C%D1%8E-Go"

	html, err := geturl.GetHtml(u)
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
