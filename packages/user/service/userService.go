package userService

import (
	"github.com/WhereNext-co/WhereNext-Backend.git/database"
	userRepo "github.com/WhereNext-co/WhereNext-Backend.git/packages/user/repo"
)

type UserServiceInterface interface {
	CreateUser(email string, password string, telNo string)
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

func (u *userService) FindUser(email string) database.UserAuth {
	return u.userRepo.FindUser(email)
}