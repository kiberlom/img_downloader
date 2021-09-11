package db

import (
	"fmt"
	"log"
	"time"

	"github.com/kiberlom/img_downloader/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB interface {
	Ping() (bool, error)
	FreeUrl() (*Url, error)
	FreeUrlAll() (*[]Url, error)
	AddNewUrl(string) error
	UpdateUrlVisit(*UpdateUrl) error
	FindUrl(string) (bool, error)
}

func createDSN(user, pass, ip, port, db string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, ip, port, db)
}

func NewConnect(cnf *config.Config) (DB, error) {
	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := createDSN(
		cnf.C.GetString("MYSQL_USER"),
		cnf.C.GetString("MYSQL_USER_PASSWORD"),
		cnf.C.GetString("MYSQL_IP"),
		cnf.C.GetString("MYSQL_PORT"),
		cnf.C.GetString("MYSQL_DATABASE"),
	)

	for i := 0; i < 100; i++ {

		time.Sleep(3 * time.Second)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Println("ERROR: не удачное соединение с MySql")
			continue
		}

		log.Println("OK: соединение с MySql установленно")
		return &con{con: db}, nil
	}

	return nil, fmt.Errorf("Не возможно соединиться с бд")

}
