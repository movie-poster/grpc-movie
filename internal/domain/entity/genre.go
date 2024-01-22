package entity

import (
	"time"

	"gorm.io/gorm"
)

type Genre struct {
	Model
	Name   string  `gorm:"column:name;type:varchar(255);not null"`
	Movies []Movie `gorm:"many2many:movie_genres;"`
	State  bool    `gorm:"column:state;type:tinyint(1);not null"`
}

func (m Genre) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return
}

func (m Genre) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now()
	return
}
