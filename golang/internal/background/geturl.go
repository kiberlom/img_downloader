package background

import "github.com/kiberlom/img_downloader/internal/db"

func SpiderUrl(con db.DB) {

	go func() {

		// сервис
		// - запрашивает новый юрл из бд
		// - соединяется с сервером по url
		// - обновляет данные по юрл в бд
		// - парсит тело юрл на наличие новых юрл
		// - добавляет данные о новых юрл в БД

	}()

}
