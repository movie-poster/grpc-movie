package entity

import (
	"time"

	"gorm.io/gorm"
)

type Director struct {
	Model
	Name      string    `gorm:"column:name;type:varchar(255);not null"`
	Birthdate time.Time `gorm:"column:birthdate;type:datetime;not null;"`
	Avatar    string    `gorm:"column:avatar;type:varchar(255);not null"`
	Movies    []Movie
	State     bool `gorm:"column:state;type:tinyint(1);not null"`
}

func (m Director) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return
}

func (m Director) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now()
	return
}
