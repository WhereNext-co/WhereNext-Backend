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

func (u *userRepo) AddFriend(user database.UserAuth, friend database.UserAuth) {
	u.dbConn.Create(&database.Friend{UserID: user.ID, FriendID: friend.ID})
	u.dbConn.Create(&database.Friend{UserID: friend.ID, FriendID: user.ID})
}

func (u *userRepo) RemoveFriend(user database.UserAuth, friend database.UserAuth) {
	u.dbConn.Delete(&database.Friend{}, "user_id = ? AND friend_id = ?", user.ID, friend.ID)
	u.dbConn.Delete(&database.Friend{}, "user_id = ? AND friend_id = ?", friend.ID, user.ID)
}

func (u *userRepo) GetFriends(user database.UserAuth) []database.UserAuth {
	var friends []database.UserAuth
	u.dbConn.Model(&user).Association("Friends").Find(&friends)
	return friends
}
