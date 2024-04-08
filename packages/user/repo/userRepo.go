package userRepo

import (
	"errors"
	"time"

	"github.com/WhereNext-co/WhereNext-Backend.git/database"
	"gorm.io/gorm"
)

type UserRepoInterface interface {

	CreateUserInfo(userName string, email string, title string, name string, 
		birthdate time.Time, region string, telNo string, profilePicture string, bio string) error
	CheckUserName(userName string) (bool, error)
}

type userRepo struct {
	dbConn *gorm.DB
}

func NewUserRepo(dbConn *gorm.DB) *userRepo {
	return &userRepo{dbConn: dbConn}
}

func (u *userRepo) CreateUserInfo(userName string, email string, title string, 
	name string, birthdate time.Time, region string, telNo string, profilePicture string, bio string) error {
    result := u.dbConn.Create(&database.User{
        UserName:       userName,
        Email:          email,
        Title:          title,
        Name:           name,
        Birthdate:      birthdate,
        Region:         region,
        TelNo:          telNo,
        ProfilePicture: profilePicture,
        Bio:            bio,
    })
    if result.Error != nil {
        return result.Error
    }
    return nil
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
func (u *userRepo) CheckUserName(userName string) (bool, error) {
    var user database.User
    result := u.dbConn.Where("user_name = ?", userName).First(&user)
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            return false, nil
        }
        return false, result.Error
    }
    return true, nil
}
