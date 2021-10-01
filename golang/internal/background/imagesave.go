package background

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/kiberlom/img_downloader/internal/db"
)

var ContentTypeImage = []string{
	"image/jpeg",
	"image/gif",
	"image/webp",
	"image/png",
}

type ImageSaveConfig struct {
	DB       db.DB
	Shutdown context.Context
	WG       *sync.WaitGroup
}

// сохранение изображений
func ImageSave(conf ImageSaveConfig) {

	defer conf.WG.Done()

	chUrl := make(chan string, 20)
	wg := &sync.WaitGroup{}

	wg.Add(2)

	go getNewUrlImage(chUrl, conf.Shutdown, wg)
	go getImgAndSave(chUrl, conf.Shutdown, wg)

}

func getNewUrlImage(chUrl chan<- string, shutdown context.Context, wg *sync.WaitGroup) {

	for {
		// получаем данные из бд адресах с изображением

	}

}

func getImgAndSave(chUrl <-chan string, shutdown context.Context, wg *sync.WaitGroup) {
	for {
		select {
		//case url <- chUrl:
		// качаем изображение
		case <-shutdown.Done():
			return
		}
	}
}

func getImg() {

	client := http.Client{
		Transport: &http.Transport{
			ResponseHeaderTimeout: 5 * time.Second,
		},
	}

	req, err := http.NewRequest("GET", "https://tn.fishki.net/26/upload/post/2021/09/17/3940023/7.jpeg", nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(len(b))

	defer res.Body.Close()

	var pathFile string

	switch http.DetectContentType(b) {
	// case "image/jpeg":

	// 	fn := getMD("33")
	// 	dir := path.Join(fn[:2], fn[2:5])

	// 	if err := os.MkdirAll(dir, 0700); err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}

	// 	pathFile = path.Join(dir, fmt.Sprintf("%s.%s", fn, "jpg"))
	// 	fmt.Println(pathFile)

	default:
		log.Println("Неизвестный файл")
		return

	}

	file, err := os.Create(pathFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	nb, err := file.Write(b)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(nb)

	fmt.Println("УСПЕХ")

}
