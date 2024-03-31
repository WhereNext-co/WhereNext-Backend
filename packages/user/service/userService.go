package userService

import (
	"log"
	"time"

	"github.com/WhereNext-co/WhereNext-Backend.git/database"
	userRepo "github.com/WhereNext-co/WhereNext-Backend.git/packages/user/repo"
)

type UserServiceInterface interface {
	CreateUser(email string, password string, telNo string)
	CreateUserInfo(title string, name string, birthdate string, region string, telNo string, profilePicture string) error
	FindUser(email string) database.UserAuth
}

type userService struct {
	userRepo userRepo.UserRepoInterface
}

func NewUserService(userRepo userRepo.UserRepoInterface) *userService {
	return &userService{userRepo}
}

func (u *userService) CreateUser(email string, password string, telNo string) {
	u.userRepo.CreateUser(email, password, telNo)
}

func (s *userService) CreateUserInfo(title string, name string, birthdate string, region string, telNo string, profilePicture string) error {
    parsedBirthdate, err := time.Parse("2006-01-02", birthdate)
    if err != nil {
        log.Printf("Error parsing birthdate: %v", err)
        return err
    }
    err = s.userRepo.CreateUserInfo(title, name, parsedBirthdate, region, telNo, profilePicture)
    if err != nil {
        return err
    }
    return nil
}

func (u *userService) FindUser(email string) database.UserAuth {
	return u.userRepo.FindUser(email)
}