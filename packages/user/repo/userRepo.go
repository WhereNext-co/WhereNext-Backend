package userRepo

import (
	"time"

	"github.com/WhereNext-co/WhereNext-Backend.git/database"
	"gorm.io/gorm"
)

type UserRepoInterface interface {
	CreateUser(email string, password string, telNo string)
	CreateUserInfo(title string, name string, birthdate time.Time, region string, telNo string, profilePicture string) error
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

func (u *userRepo) CreateUserInfo(title string, name string, birthdate time.Time, region string, telNo string, profilePicture string) error {
    result := u.dbConn.Create(&database.User{Title: title, Name: name, Birthdate: birthdate, Region: region, TelNo: telNo, ProfilePicture: profilePicture})
    if result.Error != nil {
        return result.Error
    }
    return nil
}

func (u *userRepo) FindUser(email string) database.UserAuth {
	var user database.UserAuth
	u.dbConn.First(&user, "email = ?", email)
	return user
}