package entity

import (
	"time"

	"gorm.io/gorm"
)

type Actor struct {
	Model
	Name      string    `gorm:"column:name;type:varchar(255);not null"`
	Birthdate time.Time `gorm:"column:birthdate;type:datetime;not null;" json:"birthdate"`
	Avatar    string    `gorm:"column:avatar;type:varchar(255)"`
	Movies    []Movie   `gorm:"many2many:movie_actors;"`
	State     bool      `gorm:"column:state;type:tinyint(1);not null"`
}

func (m Actor) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return
}

func (m Actor) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now()
	return
}
