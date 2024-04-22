package userService

import (
	"log"
	"time"

	"github.com/WhereNext-co/WhereNext-Backend.git/database"
	userRepo "github.com/WhereNext-co/WhereNext-Backend.git/packages/user/repo"
)

type UserServiceInterface interface {
	CreateUserInfo(Uid string, userName string, email string, title string, name string, birthdate string,
		region string, telNo string, profilePicture string, bio string) error
	CheckUserName(userName string) (bool, error)
	CheckTelephoneNumber(telNo string) (bool, error)
	FindUser(userName string) (database.User, error)
	FindUserByUid(Uid string) (database.User, error)
	UpdateUserInfo(Uid string, userName string, email string, title string, name string, birthdate string,
		region string, telNo string, profilePicture string, bio string) error
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

type userService struct {
	userRepo userRepo.UserRepoInterface
}

func NewUserService(userRepo userRepo.UserRepoInterface) *userService {
	return &userService{userRepo}
}

func (s *userService) CreateUserInfo(Uid string, userName string, email string,
	title string, name string, birthdate string, region string,
	telNo string, profilePicture string, bio string) error {
	parsedBirthdate, err := time.Parse(time.RFC3339, birthdate)
	if err != nil {
		log.Printf("Error parsing birthdate: %v", err)
		return err
	}
	err = s.userRepo.CreateUserInfo(Uid, userName, email, title, name, parsedBirthdate, region, telNo, profilePicture, bio)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) CheckUserName(userName string) (bool, error) {
	exists, err := s.userRepo.CheckUserName(userName)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *userService) CheckTelephoneNumber(telNo string) (bool, error) {
	exists, err := s.userRepo.CheckTelephoneNumber(telNo)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *userService) FindUser(userName string) (database.User, error) {
	user, err := s.userRepo.FindUser(userName)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (s *userService) FindUserByUid(Uid string) (database.User, error) {
	user, err := s.userRepo.FindUserByUid(Uid)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (s *userService) UpdateUserInfo(Uid string, userName string, email string, title string, name string, birthdate string, region string, telNo string, profilePicture string, bio string) error {
	parsedBirthdate, err := time.Parse(time.RFC3339, birthdate)
	if err != nil {
		log.Printf("Error parsing birthdate: %v", err)
		return err
	}
	err = s.userRepo.UpdateUserInfo(Uid, userName, email, title, name, parsedBirthdate, region, telNo, profilePicture, bio)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) IsFriend(UserUid string, friendName string) (bool, error) {
	isFriend, err := s.userRepo.IsFriend(UserUid, friendName)
	if err != nil {
		return false, err
	}
	return isFriend, nil
}

func (s *userService) FriendStatus(UserUid string, friendName string) (string, error) {
	friendStatus, err := s.userRepo.FriendStatus(UserUid, friendName)
	if err != nil {
		return friendStatus, err
	}
	return friendStatus, nil
}

func (s *userService) CreateFriendRequest(userUid string, friendName string) error {
	err := s.userRepo.CreateFriendRequest(userUid, friendName)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) AcceptFriendRequest(userUid string, friendName string) error {
	err := s.userRepo.AcceptFriendRequest(userUid, friendName)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) DeclineFriendRequest(userUid string, friendName string) error {
	err := s.userRepo.DeclineFriendRequest(userUid, friendName)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) CancelFriendRequest(userUid string, friendName string) error {
	err := s.userRepo.CancelFriendRequest(userUid, friendName)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) RemoveFriend(userUid string, friendName string) error {
	err := s.userRepo.RemoveFriend(userUid, friendName)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) FriendList(userUid string) ([]database.User, error) {
	friendList, err := s.userRepo.FriendList(userUid)
	if err != nil {
		return nil, err
	}
	return friendList, nil
}

func (s *userService) RequestsSent(userUid string) ([]database.FriendRequest, error) {
	requestsSent, err := s.userRepo.RequestsSent(userUid)
	if err != nil {
		return nil, err
	}
	return requestsSent, nil
}

func (s *userService) RequestsReceived(userUid string) ([]database.FriendRequest, error) {
	requestsReceived, err := s.userRepo.RequestsReceived(userUid)
	if err != nil {
		return nil, err
	}
	return requestsReceived, nil
}
