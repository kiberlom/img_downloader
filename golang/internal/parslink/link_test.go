package parslink

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
)

func BenchmarkParsUrls(b *testing.B) {

	r, err := http.Get("https://fishki.net")
	if err != nil {
		log.Fatal("Не удалось получить содержимое для парсинга ")
		return
	}

	buf := make([]byte, 100)
	var result []byte

	for {
		_, err := r.Body.Read(buf)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal("ERRROR")
			return
		}

		result = append(result, buf...)
	}

	re := string(result)

	for i := 0; i < b.N; i++ {
		s, _ := ParsAllUrlInPage(re)
		fmt.Print(len(s), "|")
	}

}
