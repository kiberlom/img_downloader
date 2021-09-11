package background

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/kiberlom/img_downloader/internal/db"
	"github.com/kiberlom/img_downloader/internal/geturl"
	"github.com/kiberlom/img_downloader/internal/parslink"
)

type ConfSpider struct {
	Shd context.Context
	WG  *sync.WaitGroup
	Con db.DB
}

func spider(con db.DB, url newUrlVisit) error {

	// получем не обработанную страницу

	fmt.Println("Обрабатывается: ", url.Url)

	// получим результат запроса
	html, err := geturl.GetHtml(url.Url)
	if err != nil {
		log.Printf("Ошибка получения html (%s): %v\n", url.Url, err)
		// обновляем статус страницы запроса
		if err := con.UpdateUrlVisit(&db.UpdateUrl{ID: url.IDDB, VisitDate: time.Now().UTC().Add(time.Hour * time.Duration(4))}); err != nil {
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
			// log.Println("Такой url уже есть в БД: ", v)
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
		ID:           url.IDDB,
		ContentType:  html.ContentType,
		CodeResponse: html.CodeRequest,
		VisitDate:    time.Now().UTC().Add(time.Hour * time.Duration(4)),
	}); err != nil {
		log.Println("Не удалось обновить данные об обработанной HTML страницы: ", err)
	}

	fmt.Printf("завершенно: найденно Всего: %d Уникальных: %d Дубликатов: %d\n", duble+unique, unique, duble)
	return nil

}

func startSpider(ctx context.Context, urls chan newUrlVisit, wg *sync.WaitGroup, con db.DB) {

	defer wg.Done()

	for {

		select {
		case <-ctx.Done():
			return
		case <-time.After(2 * time.Second):
			u, ok := <-urls
			if ok {
				if err := spider(con, u); err != nil {
					log.Println(err)
				}
				break
			}
			return
		}
	}
}

func SpiderUrl(c *ConfSpider) {
	defer c.WG.Done()

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-c.Shd.Done()
		cancel()
	}()

	wg.Add(1)
	urls := newUrlNotVisit(ctx, c.Con, wg)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go startSpider(ctx, urls, wg, c.Con)
	}

	wg.Wait()

	if err := c.Con.ResetStatisInProgress(); err != nil {
		fmt.Println("Ошибка, не удалось сбросить статус inProgress: ", err)
	}

	log.Println("Работа сервиса SpiderUrl завершена")
}
