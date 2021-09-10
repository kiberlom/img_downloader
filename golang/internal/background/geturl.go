package background

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/kiberlom/img_downloader/internal/db"
	"github.com/kiberlom/img_downloader/internal/geturl"
	"github.com/kiberlom/img_downloader/internal/parslink"
)

func SpiderUrl(con db.DB, wg *sync.WaitGroup) {

	go func() {

		defer wg.Done()
		// сервис
		// - запрашивает новый юрл из бд
		// - соединяется с сервером по url
		// - обновляет данные по юрл в бд
		// - парсит тело юрл на наличие новых юрл
		// - добавляет данные о новых юрл в БД

		for {
			time.Sleep(2 * time.Second)

			// получем не обработанную страницу
			url, err := con.FreeUrl()
			if err != nil {
				log.Println("Ошибка получения URL из бд")
				continue
			}
			fmt.Println("Обрабатывается: ", url.Url)

			// получим результат запроса
			html, err := geturl.GetHtml(url.Url)
			if err != nil {
				log.Printf("Ошибка получения html (%s): %v\n", url.Url, err)
				continue
			}

			// парсим страницу на url
			links, err := parslink.GetAllUrl(string(html.Body))
			if err != nil {
				log.Println("Ощибка парсинга теля html страницы")
			}

			// добавляем новый url в базу
			for _, v := range links {
				if err := con.AddNewUrl(v); err != nil {
					log.Println("Ошибка в БД не добавленн: ", v)
				}

			}

			// обновляем статус страницы запроса
			if err := con.UpdateUrlVisit(&db.UpdateUrl{
				ID:           url.ID,
				ContentType:  html.ContentType,
				CodeResponse: html.CodeRequest,
				VisitDate:    time.Now().UTC().Add(time.Hour * time.Duration(4)),
			}); err != nil {
				log.Panicln("Не удалось обновить данные об обработанной HTML страницы: ", err)
			}

			fmt.Println("завершенно")

		}

	}()

}
