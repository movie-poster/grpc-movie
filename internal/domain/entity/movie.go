package entity

import (
	"time"

	"gorm.io/gorm"
)

type Movie struct {
	Model
	Title      string   `gorm:"column:title;type:varchar(255);not null"`
	Synopsis   string   `gorm:"column:synopsis;type:varchar(1000);not null"`
	Year       uint32   `gorm:"column:year;type:int(11) unsigned;not null"`
	Rating     float64  `gorm:"column:rating;type:decimal(10,2);not null"`
	Duration   uint32   `gorm:"column:duration;type:int(11) unsigned;not null"`
	DirectorID uint64   `gorm:"column:director_id;type:bigint(20) unsigned;not null"`
	Director   Director `gorm:"foreignKey:DirectorID"`
	Actors     []Actor  `gorm:"many2many:movie_actors;"`
	Genres     []Genre  `gorm:"many2many:movie_genres;"`
	State      bool     `gorm:"column:state;type:tinyint(1);not null"`
}

func (m Movie) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return
}

func (m Movie) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = time.Now()
	return
}
