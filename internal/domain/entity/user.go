package entity

type User struct {
	ID       uint64 `gorm:"column:id;primary_key;auto_increment;"`
	Name     string `gorm:"column:name;not null;"`
	Document string `gorm:"column:document;not null;"`
	Phone    string `gorm:"column:phone;not null;"`
}
