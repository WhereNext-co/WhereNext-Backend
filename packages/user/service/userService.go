package userService

import (
	"log"
	"time"

	"github.com/WhereNext-co/WhereNext-Backend.git/database"
	userRepo "github.com/WhereNext-co/WhereNext-Backend.git/packages/user/repo"
)

type UserServiceInterface interface {
	CreateUserInfo(userName string, email string, title string, name string, birthdate string,
		region string, telNo string, profilePicture string, bio string) error
	CheckUserName(userName string) (bool, error)
	FindUser(userName string) (database.User, error)
	UpdateUserInfo(userName string, email string, title string, name string, birthdate string,
		region string, telNo string, profilePicture string, bio string) error
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

type userService struct {
	userRepo userRepo.UserRepoInterface
}

func NewUserService(userRepo userRepo.UserRepoInterface) *userService {
	return &userService{userRepo}
}

func (s *userService) CreateUserInfo(userName string, email string,
	title string, name string, birthdate string, region string,
	telNo string, profilePicture string, bio string) error {
	parsedBirthdate, err := time.Parse("2006-01-02", birthdate)
	if err != nil {
		log.Printf("Error parsing birthdate: %v", err)
		return err
	}
	err = s.userRepo.CreateUserInfo(userName, email, title, name, parsedBirthdate, region, telNo, profilePicture, bio)
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

func (s *userService) FindUser(userName string) (database.User, error) {
	user, err := s.userRepo.FindUser(userName)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (s *userService) UpdateUserInfo(userName string, email string, title string, name string, birthdate string, region string, telNo string, profilePicture string, bio string) error {
	parsedBirthdate, err := time.Parse("2006-01-02", birthdate)
	if err != nil {
		log.Printf("Error parsing birthdate: %v", err)
		return err
	}
	err = s.userRepo.UpdateUserInfo(userName, email, title, name, parsedBirthdate, region, telNo, profilePicture, bio)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) IsFriend(userName string, friendName string) (bool, error) {
	isFriend, err := s.userRepo.IsFriend(userName, friendName)
	if err != nil {
		return false, err
	}
	return isFriend, nil
}

func (s *userService) CreateFriendRequest(userName string, friendName string) error {
	err := s.userRepo.CreateFriendRequest(userName, friendName)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) AcceptFriendRequest(userName string, friendName string) error {
	err := s.userRepo.AcceptFriendRequest(userName, friendName)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) DeclineFriendRequest(userName string, friendName string) error {
	err := s.userRepo.DeclineFriendRequest(userName, friendName)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) CancelFriendRequest(userName string, friendName string) error {
	err := s.userRepo.CancelFriendRequest(userName, friendName)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) RemoveFriend(userName string, friendName string) error {
	err := s.userRepo.RemoveFriend(userName, friendName)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) FriendList(userName string) ([]database.User, error) {
	friendList, err := s.userRepo.FriendList(userName)
	if err != nil {
		return nil, err
	}
	return friendList, nil
}

func (s *userService) RequestsSent(userName string) ([]database.FriendRequest, error) {
	requestsSent, err := s.userRepo.RequestsSent(userName)
	if err != nil {
		return nil, err
	}
	return requestsSent, nil
}

func (s *userService) RequestsReceived(userName string) ([]database.FriendRequest, error) {
	requestsReceived, err := s.userRepo.RequestsReceived(userName)
	if err != nil {
		return nil, err
	}
	return requestsReceived, nil
}
