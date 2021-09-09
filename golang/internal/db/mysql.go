package db

import "gorm.io/gorm"

type con struct {
	con *gorm.DB
}

type Url struct {
	ID          int    `gorm:"id"`
	Url         string `gorm:"url"`
	ContentType string `gorm:"content_type"`
}

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

func (c *con) FreeUrl() (*Url, error) {
	r := new(Url)
	tx := c.con.Table("url").Take(r)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return r, nil
}

func (c *con) AddNewUrl(url string) error {
	tx := c.con.Table("url").Create(&Url{Url: "jkfkjdskfjkdsjfkj", ContentType: "image/jpeg"})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
