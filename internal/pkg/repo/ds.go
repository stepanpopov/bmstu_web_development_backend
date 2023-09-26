package repo

import (
	"gorm.io/gorm"
)

type Product struct {
	ID        uint `gorm:"primarykey"`
	Code  string
	Price uint
}
