package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/kiberlom/img_downloader/internal/geturl"
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

	//u := "https://img5.goodfon.ru"
	u := "https://play.google.com/store/apps/details?id=com.goodfon.goodfon&rdid=com.goodfon.goodfon"

	html, err := geturl.GetHtml(u)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(html)

	rxg, err := regexp.Compile(`(http|ftp|https)://([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?`)
	if err != nil {
		log.Fatal(err)
	}

	sr := rxg.FindAllString(html, -1)

	f := make(map[string]struct{})
	for _, v := range sr {
		if _, ok := f[v]; !ok {
			f[v] = struct{}{}
		}
	}
	for i := range f {
		fmt.Println(i)
	}

}
