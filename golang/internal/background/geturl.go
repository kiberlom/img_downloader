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

func start(con db.DB) error {

	// получем не обработанную страницу
	url, err := con.FreeUrl()
	if err != nil {
		return fmt.Errorf("Ошибка получения URL из бд")
	}
	fmt.Println("Обрабатывается: ", url.Url)

	// получим результат запроса
	html, err := geturl.GetHtml(url.Url)
	if err != nil {
		log.Printf("Ошибка получения html (%s): %v\n", url.Url, err)
		// обновляем статус страницы запроса
		if err := con.UpdateUrlVisit(&db.UpdateUrl{ID: url.ID, VisitDate: time.Now().UTC().Add(time.Hour * time.Duration(4))}); err != nil {
			log.Panicln("Не удалось обновить данные об обработанной HTML страницы: ", err)
		}
		return fmt.Errorf("Ошибка получения html (%s): %v\n", url.Url, err)
	}

	// парсим страницу на url
	links, err := parslink.GetAllUrl(string(html.Body))
	if err != nil {
		return fmt.Errorf("Ощибка парсинга теля html страницы")
	}

	var duble int
	var unique int
	// добавляем новый url в базу
	for _, v := range links {

		// проверим url в БД
		ex, err := con.FindUrl(v)
		if err != nil {
			log.Println("Ошибка проверки в БД повторный url: ", v)
		}

		if ex {
			log.Println("Такой url уже есть в БД: ", v)
			duble++
			continue
		}

		if err := con.AddNewUrl(v); err != nil {
			log.Println("Ошибка в БД не добавленн: ", v)
		}

		unique++

	}

	// обновляем статус страницы запроса
	if err := con.UpdateUrlVisit(&db.UpdateUrl{
		ID:           url.ID,
		ContentType:  html.ContentType,
		CodeResponse: html.CodeRequest,
		VisitDate:    time.Now().UTC().Add(time.Hour * time.Duration(4)),
	}); err != nil {
		log.Println("Не удалось обновить данные об обработанной HTML страницы: ", err)
	}

	fmt.Printf("завершенно: найденно Всего: %d Уникальных: %d Дубликатов: %d\n", duble+unique, unique, duble)
	return nil

}

func SpiderUrl(con db.DB, wg *sync.WaitGroup) {

	defer wg.Done()
	for {
		time.Sleep(2 * time.Second)
		start(con)
	}

}
