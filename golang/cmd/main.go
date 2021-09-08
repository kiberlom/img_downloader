package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func getUrl() (*[]byte, error) {
	client := http.Client{}

	transport := http.Transport{ResponseHeaderTimeout: 5 * time.Second}

	client.Transport = &transport

	r, err := http.NewRequest("GET", "https://img5.goodfon.ru/original/2048x1365/4/ea/park-doroga-utro.jpg", nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36")

	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.Header.Get("Content-Type") != "image/jpeg" {
		return nil, fmt.Errorf("%s", resp.Header.Get("Content-Type"))
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &b, nil
}

func createFile(b *[]byte) error {

	f, err := os.Create("img.jpeg")
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err := f.Write(*b); err != nil {
		return err
	}

	return nil

}

func main() {

	content, err := getUrl()
	if err != nil {
		log.Fatal(err)
	}

	if err := createFile(content); err != nil {
		log.Fatal(err)
	}

}
