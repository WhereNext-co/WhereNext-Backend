package userService

import (
	"log"
	"time"

	userRepo "github.com/WhereNext-co/WhereNext-Backend.git/packages/user/repo"
)

type UserServiceInterface interface {
	CreateUserInfo(userName string, email string, title string, name string, birthdate string, 
		region string, telNo string, profilePicture string, bio string) error
	CheckUserName(userName string) (bool, error)
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