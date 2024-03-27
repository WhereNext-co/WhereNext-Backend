package userRepo

import (
	"github.com/WhereNext-co/WhereNext-Backend.git/database"
	"gorm.io/gorm"
)

type UserRepoInterface interface {
	CreateUser(email string, password string, telNo string)
	FindUser(email string) database.UserAuth
}

type userRepo struct {
	dbConn *gorm.DB
}

func NewUserRepo(dbConn *gorm.DB) *userRepo {
	return &userRepo{dbConn: dbConn}
}

func (u *userRepo) CreateUser(email string, password string, telNo string) {
	u.dbConn.Create(&database.UserAuth{Email: email, Password: password, TelNo: telNo})
}

func (u *userRepo) FindUser(email string) database.UserAuth {
	var user database.UserAuth
	u.dbConn.First(&user, "email = ?", email)
	return user
}