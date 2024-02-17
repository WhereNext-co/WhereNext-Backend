package database

import (
	"gorm.io/gorm"
)

type UserAuth struct {
	gorm.Model
	Email    string
	Password string
	TelNo    string
}
