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
	"github.com/sirupsen/logrus"
)

type ConfBackService struct {
	Shd context.Context
	WG  *sync.WaitGroup
	Con db.DB
	Log *logrus.Logger
}

func spiderGet(ctx context.Context, con db.DB, url newUrlVisit, index int) error {

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
		return fmt.Errorf("Ошибка парсинга теля html страницы")
	}

	var double int
	var unique int
	// добавляем новый url в базу

	tc := con.TransStart()

	for _, v := range links {
		select {
		case <-ctx.Done():
			tc.TransError()
			fmt.Println("Принудительное завершение")
			return nil
		default:
			// проверим url в БД
			ex, err := tc.FindUrl(v)
			if err != nil {
				log.Println("Ошибка проверки в БД повторный url: ", v)
			}

			if ex {
				// log.Println("Такой url уже есть в БД: ", v)
				double++
				continue
			}

			if err := tc.AddNewUrl(v); err != nil {
				log.Println("Ошибка в БД не добавлен: ", v)
			}

			unique++
		}

	}

	// обновляем статус страницы запроса
	if err := tc.UpdateUrlVisit(&db.UpdateUrl{
		ID:           url.IDDB,
		ContentType:  html.ContentType,
		CodeResponse: html.CodeRequest,
		VisitDate:    time.Now().UTC().Add(time.Hour * time.Duration(4)),
	}); err != nil {
		log.Println("Не удалось обновить данные об обработанной HTML страницы: ", err)
	}

	tc.TransCommit()
	fmt.Println("--------------[ ", index, " ]-----------------")
	fmt.Println("Обрабатывается: ", url.Url)
	fmt.Printf("завершено: найдено Всего: %d Уникальных: %d Дубликатов: %d\n", double+unique, unique, double)
	return nil

}

func startSpider(ctx context.Context, i int, urls chan newUrlVisit, wg *sync.WaitGroup, con db.DB, index int) {

	defer wg.Done()

	for {

		select {
		case <-ctx.Done():
			return
		case <-time.After(2 * time.Millisecond):
			u, ok := <-urls
			if ok {
				if err := spiderGet(ctx, con, u, index); err != nil {
					log.Println(err)
				}
				break
			}
			return
		}
	}
}

func SpiderUrl(c *ConfBackService) {
	defer c.WG.Done()

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-c.Shd.Done()
		fmt.Println("Command to stop")
		cancel()
	}()

	wg.Add(1)
	urls := newUrlNotVisit(ctx, c.Con, wg)

	for i := 1; i < 3; i++ {
		wg.Add(1)
		go startSpider(ctx, i, urls, wg, c.Con, i)
	}

	wg.Wait()

	if err := c.Con.ResetStatisInProgress(); err != nil {
		fmt.Println("Ошибка, не удалось сбросить статус inProgress: ", err)
	}

	fmt.Println("Работа сервиса SpiderUrl завершена")
}
