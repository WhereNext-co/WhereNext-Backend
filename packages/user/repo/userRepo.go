package userRepo

import (
	"errors"
	"time"

	"github.com/WhereNext-co/WhereNext-Backend.git/database"
	"gorm.io/gorm"
)

type UserRepoInterface interface {
	CreateUserInfo(Uid string, userName string, email string, title string, name string,
		birthdate time.Time, region string, telNo string, profilePicture string, bio string) error
	CheckUserName(userName string) (bool, error)
	CheckTelephoneNumber(telNo string) (bool, error)
	FindUser(userName string) (database.User, error)
	FindUserByUid(Uid string) (database.User, error)
	UpdateUserInfo(Uid string, userName string, email string, title string, name string,
		birthdate time.Time, region string, telNo string, profilePicture string, bio string) error
	IsFriend(Uid string, friendName string) (bool, error)
	FriendStatus(Uid string, friendName string) (string, error)
	CreateFriendRequest(Uid string, friendName string) error
	AcceptFriendRequest(Uid string, friendName string) error
	DeclineFriendRequest(Uid string, friendName string) error
	CancelFriendRequest(Uid string, friendName string) error
	RemoveFriend(Uid string, friendName string) error
	FriendList(Uid string) ([]database.User, error)
	RequestsSent(Uid string) ([]database.FriendRequest, error)
	RequestsReceived(Uid string) ([]database.FriendRequest, error)
}

type userRepo struct {
	dbConn *gorm.DB
}

func NewUserRepo(dbConn *gorm.DB) *userRepo {
	return &userRepo{dbConn: dbConn}
}

func (u *userRepo) CreateUserInfo(Uid string, userName string, email string, title string,
	name string, birthdate time.Time, region string, telNo string, profilePicture string, bio string) error {
	result := u.dbConn.Create(&database.User{
		Uid:            Uid,
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

func (u *userRepo) CheckTelephoneNumber(telNo string) (bool, error) {
    var user database.User
    result := u.dbConn.Where("tel_no = ?", telNo).First(&user)
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

func (u *userRepo) FindUserByUid(Uid string) (database.User, error) {
	var user database.User
	result := u.dbConn.Where("uid = ?", Uid).First(&user)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func (u *userRepo) UpdateUserInfo(Uid string, userName string, email string, title string,
	name string, birthdate time.Time, region string, telNo string, profilePicture string, bio string) error {
	user, err := u.FindUserByUid(Uid)
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

// IsFriend need fixing, always return false
func (u *userRepo) IsFriend(userUid string, friendName string) (bool, error) {
	user, err := u.FindUserByUid(userUid)
	if err != nil {
		return false, err
	}
	if user.UserName == friendName {
		return false, errors.New("cannot be friend with yourself")
	}
	friend, err := u.FindUser(friendName)
	if err != nil {
		return false, err
	}
	var userFriend []database.User
	err = u.dbConn.Model(&user).Association("Friends").Find(&userFriend)
	if err != nil {
		return false, err
	}
	for _, f := range userFriend {
		if f.Uid == friend.Uid {
			return true, nil
		}
	}
	return false, nil
}

func (u *userRepo) FriendStatus(userUid string, friendName string) (string, error) {
	user, err := u.FindUserByUid(userUid)
	if err != nil {
		return "cannot find user", err
	}
	if user.UserName == friendName {
		return "cannot be friend with yourself", errors.New("cannot be friend with yourself")
	}
	friend, err := u.FindUser(friendName)
	if err != nil {
		return "cannot find friend", err
	}
	isFriend, err := u.IsFriend(userUid, friendName)
	if err != nil {
		return "error checking friend status", err
	}
	if isFriend {
		return "already friends", nil
	}
	var request database.FriendRequest
	result := u.dbConn.Where("sender_uid = ? AND receiver_uid = ? AND status = ?", user.Uid, friend.Uid, "Pending").First(&request)
	if result.Error == nil {
		return "friend request already sent", nil
	}
	result = u.dbConn.Where("sender_uid = ? AND receiver_uid = ? AND status = ?", friend.Uid, user.Uid, "Pending").First(&request)
	if result.Error == nil {
		return "friend request received", nil
	}
	return "not friends", nil
}

func (u *userRepo) CreateFriendRequest(userUid string, friendName string) error {
	user, err := u.FindUserByUid(userUid)
	if err != nil {
		return err
	}
	if user.UserName == friendName {
		return errors.New("cannot send friend request to yourself")
	}
	friend, err := u.FindUser(friendName)
	if err != nil {
		return err
	}
	isFriend, err := u.IsFriend(userUid, friendName)
	if err != nil {
		return err
	}
	if isFriend {
		return errors.New("already friends")
	}
	var request database.FriendRequest
	isSent := u.dbConn.Where("sender_id = ? AND receiver_id = ? AND status = ?", user.Uid, friend.Uid, "Pending").First(&request)
	if isSent.Error == nil {
		return errors.New("friend request already sent")
	}
	result := u.dbConn.Create(&database.FriendRequest{
		SenderUid:   user.Uid,
		ReceiverUid: friend.Uid,
		Status:      "Pending",
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *userRepo) AcceptFriendRequest(userUid string, friendName string) error {
	user, err := u.FindUserByUid(userUid)
	if err != nil {
		return err
	}
	if user.UserName == friendName {
		return errors.New("cannot send friend request to yourself")
	}
	friend, err := u.FindUser(friendName)
	if err != nil {
		return err
	}
	isFriend, err := u.IsFriend(userUid, friendName)
	if err != nil {
		return err
	}
	if isFriend {
		return errors.New("already friends")
	}
	var request database.FriendRequest
	result := u.dbConn.Where("sender_uid = ? AND receiver_uid = ? AND status = ?", friend.Uid, user.Uid, "Pending").First(&request)
	if result.Error != nil {
		return result.Error
	}
	result = u.dbConn.Model(&request).Update("status", "Accepted")
	if result.Error != nil {
		return result.Error
	}
	u.dbConn.Where("uid = ?", userUid).Preload("Friends").First(&user)
	u.dbConn.Where("user_name = ?", friendName).Preload("Friends").First(&friend)
	u.dbConn.Model(&user).Association("Friends").Append(&friend)
	u.dbConn.Model(&friend).Association("Friends").Append(&user)

	return nil
}

func (u *userRepo) DeclineFriendRequest(userUid string, friendName string) error {

	user, err := u.FindUserByUid(userUid)
	if err != nil {
		return err
	}
	if user.UserName == friendName {
		return errors.New("cannot send friend request to yourself")
	}
	friend, err := u.FindUser(friendName)
	if err != nil {
		return err
	}
	isFriend, err := u.IsFriend(userUid, friendName)
	if err != nil {
		return err
	}
	if isFriend {
		return errors.New("already friends")
	}
	var request database.FriendRequest
	result := u.dbConn.Where("sender_uid = ? AND receiver_uid = ? AND status = ?", friend.Uid, user.Uid, "Pending").First(&request)
	if result.Error != nil {
		return result.Error
	}
	result = u.dbConn.Model(&request).Update("status", "Declined")
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *userRepo) CancelFriendRequest(userUid string, friendName string) error {
	user, err := u.FindUserByUid(userUid)
	if err != nil {
		return err
	}
	if user.UserName == friendName {
		return errors.New("cannot send friend request to yourself")
	}
	friend, err := u.FindUser(friendName)
	if err != nil {
		return err
	}
	isFriend, err := u.IsFriend(userUid, friendName)
	if err != nil {
		return err
	}
	if isFriend {
		return errors.New("already friends")
	}
	var request database.FriendRequest
	result := u.dbConn.Where("sender_uid = ? AND receiver_uid = ? AND status = ?", user.Uid, friend.Uid, "Pending").First(&request)
	if result.Error != nil {
		return result.Error
	}
	result = u.dbConn.Model(&request).Update("status", "Declined")
	if result.Error != nil {
		return result.Error
	}
	return nil
}

type Friendship struct {
	UserID   string
	FriendID string
}

func (u *userRepo) RemoveFriend(userUid string, friendName string) error {
	user, err := u.FindUserByUid(userUid)
	if err != nil {
		return err
	}
	friend, err := u.FindUser(friendName)
	if err != nil {
		return err
	}
	result := u.dbConn.Exec("DELETE FROM user_profiles WHERE user_uid = ? AND friend_uid = ?", user.Uid, friend.Uid)
	if result.Error != nil {
		return result.Error
	}
	result = u.dbConn.Exec("DELETE FROM user_profiles WHERE user_uid = ? AND friend_uid = ?", friend.Uid, user.Uid)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *userRepo) FriendList(userUid string) ([]database.User, error) {
	var user database.User
	result := u.dbConn.Where("uid = ?", userUid).Preload("Friends").First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user.Friends, nil
}

func (u *userRepo) RequestsSent(userUid string) ([]database.FriendRequest, error) {
	var user database.User
	result := u.dbConn.Where("uid = ?", userUid).Preload("RequestsSent").First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user.RequestsSent, nil
}

func (u *userRepo) RequestsReceived(userUid string) ([]database.FriendRequest, error) {
	var user database.User
	result := u.dbConn.Where("uid = ?", userUid).Preload("RequestsReceived.Sender").Preload("RequestsReceived", func(db *gorm.DB) *gorm.DB {
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
