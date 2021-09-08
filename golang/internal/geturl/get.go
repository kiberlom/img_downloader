package geturl

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type req struct {
	Body        []byte
	CodeRequest int
	ContentType string
}

func clientGet(url string) (*req, error) {
	client := http.Client{}

	transport := http.Transport{ResponseHeaderTimeout: 5 * time.Second}

	client.Transport = &transport

	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36")

	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &req{
		Body:        b,
		ContentType: resp.Header.Get("Content-type"),
		CodeRequest: resp.StatusCode}, nil
}

// закачка изображения
func GetImage(url string) (*[]byte, error) {

	b, err := clientGet(url)
	if err != nil {
		return nil, err
	}

	switch b.ContentType {
	case "image/png", "image/jpeg", "image/webp":
		return &b.Body, nil
	default:
		return nil, fmt.Errorf("неизвестный тип")
	}

}

// получение html страницы
func GetHtml(url string) (string, error) {

	r, err := clientGet(url)
	if err != nil {
		return "", err
	}

	if r.CodeRequest == http.StatusOK {
		s := string(r.Body)
		return s, nil
	}

	return "", fmt.Errorf("err status request: %d", r.CodeRequest)

}
