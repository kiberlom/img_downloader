package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type con struct {
	con *gorm.DB
}

type Url struct {
	ID           int     `gorm:"id"`
	Url          string  `gorm:"url"`
	ContentType  string  `gorm:"content_type"`
	CodeResponse *int    `gorm:"code_response"`
	Visit        *string `gorm:"visit"`
	Status       *string `gorm:"status"`
}

type UpdateUrl struct {
	ID           int
	ContentType  string
	CodeResponse int
	VisitDate    time.Time
}

// транзакции
func (c *con) TransStart() DB {
	return &con{
		con: c.con.Begin(),
	}
}

func (c *con) TransError() {
	c.con.Rollback()
}

func (c *con) TransCommit() {
	c.con.Commit()
}

// методы
func (c *con) Ping() (bool, error) {
	sql, err := c.con.DB()
	if err != nil {
		return false, err
	}

	if err := sql.Ping(); err != nil {
		return false, nil
	}

	return true, nil
}

func (c *con) FindUrl(u string) (bool, error) {
	tx := c.con.Table("url").Where("url = ?", u).Find(&Url{})
	if tx.Error != nil {
		return false, tx.Error
	}

	return tx.RowsAffected > 0, nil
}

func (c *con) FreeUrl() (*Url, error) {
	r := new(Url)
	tx := c.con.Table("url").Where("visit IS NULL AND status IS NULL").First(r)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return r, nil
}

func (c *con) FreeUrlAll() (*[]Url, error) {
	r := new([]Url)
	tx := c.con.Table("url").Where("visit IS NULL").Find(r)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return r, nil
}

func (c *con) AddNewUrl(url string) error {
	tx := c.con.Table("url").Create(&Url{Url: url})
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return fmt.Errorf("not found")
	}

	return nil
}

func (c *con) UpdateUrlVisit(u *UpdateUrl) error {

	date := u.VisitDate.Format("2006-01-02 15:04:05")
	fmt.Println(date)

	s := "visit"
	tx := c.con.Table("url").Model(&Url{}).Where("id = ?", u.ID).Updates(Url{ContentType: u.ContentType, Visit: &date, CodeResponse: &u.CodeResponse, Status: &s})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (c *con) UpdateUrlStatus(id int) error {

	tx := c.con.Table("url").Model(&Url{}).Where("id = ?", id).Update("status", "inProgres")
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (c *con) ResetStatisInProgress() error {
	tx := c.con.Table("url").Model(&Url{}).Where("status = 'inProgres'").Updates(map[string]interface{}{"status": nil})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
