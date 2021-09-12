package geturl

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Req struct {
	Body        []byte
	CodeRequest int
	ContentType string
}

func clientGet(url string) (*Req, error) {
	client := http.Client{
		// CheckRedirect: func(req *http.Request, via []*http.Request) error {
		// 	return errors.New("Redirect")
		// },
	}

	transport := http.Transport{ResponseHeaderTimeout: 2 * time.Second}

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

	return &Req{
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
func GetHtml(url string) (*Req, error) {

	r, err := clientGet(url)
	if err != nil {
		return nil, err
	}

	// s := string(r.Body)
	return &Req{Body: r.Body, ContentType: r.ContentType, CodeRequest: r.CodeRequest}, nil

}
