package webserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"text/template"

	"github.com/kiberlom/img_downloader/internal/db"
	"golang.org/x/net/websocket"
)

type WebServer struct {
	DB       db.DB
	Shutdown context.Context
	WG       *sync.WaitGroup
}

func (w *WebServer) home(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("WWWWWWWWWWWWWWWWWWWW")

	// rw.Write([]byte("HELLO"))

	tmpl, err := template.ParseFiles("templates/html/index.html")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	tmpl.Execute(rw, nil)
}

func NewWebServer(w *WebServer) {
	log.Println("WEB server start ...")
	defer w.WG.Done()

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":80",
		Handler: mux,
	}

	mux.HandleFunc("/", w.home)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("templates/static"))))
	mux.Handle("/ws", websocket.Handler(func(conn *websocket.Conn) {

	}))

	go func() {
		err := server.ListenAndServe()
		if err.Error() == "http: Server closed" {
			log.Println("WEB server stop")
			return
		}

		if err != nil {
			fmt.Println(err)
		}

	}()

	<-w.Shutdown.Done()
	if err := server.Close(); err != nil {
		log.Println(err)
	}

}
