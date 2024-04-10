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
	FindUser(userName string) (database.User, error)
	UpdateUserInfo(userName string, email string, title string, name string,
		birthdate time.Time, region string, telNo string, profilePicture string, bio string) error
	IsFriend(userName string, friendName string) (bool, error)
	CreateFriendRequest(userName string, friendName string) error
	AcceptFriendRequest(userName string, friendName string) error
	DeclineFriendRequest(userName string, friendName string) error
	CancelFriendRequest(userName string, friendName string) error
	RemoveFriend(userName string, friendName string) error
	FriendList(userName string) ([]database.User, error)
	RequestsSent(userName string) ([]database.FriendRequest, error)
	RequestsReceived(userName string) ([]database.FriendRequest, error)
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

func (u *userRepo) FindUser(userName string) (database.User, error) {
	var user database.User
	result := u.dbConn.Where("user_name = ?", userName).First(&user)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func (u *userRepo) UpdateUserInfo(userName string, email string, title string,
	name string, birthdate time.Time, region string, telNo string, profilePicture string, bio string) error {
	user, err := u.FindUser(userName)
	if err != nil {
		return err
	}
	user.Email = email
	user.Title = title
	user.Name = name
	user.Birthdate = birthdate
	user.Region = region
	user.TelNo = telNo
	user.ProfilePicture = profilePicture
	user.Bio = bio
	result := u.dbConn.Save(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *userRepo) IsFriend(userName string, friendName string) (bool, error) {
	if userName == friendName {
		return false, errors.New("Cannot be friend with yourself")
	}
	user, err := u.FindUser(userName)
	if err != nil {
		return false, err
	}
	friend, err := u.FindUser(friendName)
	if err != nil {
		return false, err
	}
	for _, f := range user.Friends {
		if f.ID == friend.ID {
			return true, nil
		}
	}
	return false, nil
}

func (u *userRepo) CreateFriendRequest(userName string, friendName string) error {
	if userName == friendName {
		return errors.New("Cannot send friend request to yourself")
	}
	user, err := u.FindUser(userName)
	if err != nil {
		return err
	}
	friend, err := u.FindUser(friendName)
	if err != nil {
		return err
	}
	isfriend, err := u.IsFriend(userName, friendName)
	if isfriend {
		return errors.New("Already friends")
	}
	var request database.FriendRequest
	issent := u.dbConn.Where("sender_id = ? AND receiver_id = ? AND status = ?", user.ID, friend.ID, "Pending").First(&request)
	if issent.Error == nil {
		return errors.New("Friend request already sent")
	}
	result := u.dbConn.Create(&database.FriendRequest{
		SenderID:   user.ID,
		ReceiverID: friend.ID,
		Status:     "Pending",
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *userRepo) AcceptFriendRequest(userName string, friendName string) error {
	if userName == friendName {
		return errors.New("Cannot send friend request to yourself")
	}
	user, err := u.FindUser(userName)
	if err != nil {
		return err
	}
	friend, err := u.FindUser(friendName)
	if err != nil {
		return err
	}
	isfriend, err := u.IsFriend(userName, friendName)
	if isfriend {
		return errors.New("Already friends")
	}
	var request database.FriendRequest
	result := u.dbConn.Where("sender_id = ? AND receiver_id = ? AND status = ?", friend.ID, user.ID, "Pending").First(&request)
	if result.Error != nil {
		return result.Error
	}
	result = u.dbConn.Model(&request).Update("status", "Accepted")
	if result.Error != nil {
		return result.Error
	}
	u.dbConn.Where("user_name = ?", userName).Preload("Friends").First(&user)
	u.dbConn.Where("user_name = ?", friendName).Preload("Friends").First(&friend)
	u.dbConn.Model(&user).Association("Friends").Append(&friend)
	u.dbConn.Model(&friend).Association("Friends").Append(&user)

	return nil
}

func (u *userRepo) DeclineFriendRequest(userName string, friendName string) error {
	if userName == friendName {
		return errors.New("Cannot send friend request to yourself")
	}
	user, err := u.FindUser(userName)
	if err != nil {
		return err
	}
	friend, err := u.FindUser(friendName)
	if err != nil {
		return err
	}
	isfriend, err := u.IsFriend(userName, friendName)
	if isfriend {
		return errors.New("Already friends")
	}
	var request database.FriendRequest
	result := u.dbConn.Where("sender_id = ? AND receiver_id = ? AND status = ?", friend.ID, user.ID, "Pending").First(&request)
	if result.Error != nil {
		return result.Error
	}
	result = u.dbConn.Model(&request).Update("status", "Declined")
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *userRepo) CancelFriendRequest(userName string, friendName string) error {
	if userName == friendName {
		return errors.New("Cannot send friend request to yourself")
	}
	user, err := u.FindUser(userName)
	if err != nil {
		return err
	}
	friend, err := u.FindUser(friendName)
	if err != nil {
		return err
	}
	isfriend, err := u.IsFriend(userName, friendName)
	if isfriend {
		return errors.New("Already friends")
	}
	var request database.FriendRequest
	result := u.dbConn.Where("sender_id = ? AND receiver_id = ? AND status = ?", user.ID, friend.ID, "Pending").First(&request)
	if result.Error != nil {
		return result.Error
	}
	result = u.dbConn.Model(&request).Update("status", "Declined")
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *userRepo) RemoveFriend(userName string, friendName string) error {
	user, err := u.FindUser(userName)
	if err != nil {
		return err
	}
	friend, err := u.FindUser(friendName)
	if err != nil {
		return err
	}
	result := u.dbConn.Model(&user).Delete(&friend)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *userRepo) FriendList(userName string) ([]database.User, error) {
	var user database.User
	result := u.dbConn.Where("user_name = ?", userName).Preload("Friends").First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user.Friends, nil
}

func (u *userRepo) RequestsSent(userName string) ([]database.FriendRequest, error) {
	var user database.User
	result := u.dbConn.Where("user_name = ?", userName).Preload("RequestsSent").First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user.RequestsSent, nil
}

func (u *userRepo) RequestsReceived(userName string) ([]database.FriendRequest, error) {
	var user database.User
	result := u.dbConn.Where("user_name = ?", userName).Preload("RequestsReceived.Sender").Preload("RequestsReceived", func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", "pending")
	}).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user.RequestsReceived, nil
}

// userRepo works left
// 1. FriendList query only the username, name, and profile picture
// 2. RequestsReceived query only the ID CreatedAt the username, name, and profile picture
