package background

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/kiberlom/img_downloader/internal/db"
)

// сервис берет данные о новых не проверенных адресах

type newUrlVisit struct {
	IDDB int
	Url  string
}

func newUrlNotVisit(ctx context.Context, con db.DB, wg *sync.WaitGroup) chan newUrlVisit {

	ch := make(chan newUrlVisit, 100)

	go func() {

		defer func() {
			close(ch)
			wg.Done()
		}()

		for {
			//time.Sleep(500 * time.Millisecond)
			select {
			case <-ctx.Done():
				return
			default:
				getNewUrl(con, ch)
			}

		}
	}()

	return ch

}

func getNewUrl(con db.DB, ch chan newUrlVisit) {

	tc := con.TransStart()
	u, err := tc.FreeUrl()
	if err != nil {
		tc.TransError()
		log.Println("Не возможно получить новые URL: ", err)
		return
	}

	if err := tc.UpdateUrlStatus(u.ID); err != nil {
		log.Println("Не возможно обновить статус: ", err)
		tc.TransError()
		return
	}

	select {
	case ch <- newUrlVisit{IDDB: u.ID, Url: u.Url}:
		tc.TransCommit()
		//fmt.Printf("%+v\n", u)
	default:
		tc.TransError()
		time.Sleep(1 * time.Second)
	}

}
