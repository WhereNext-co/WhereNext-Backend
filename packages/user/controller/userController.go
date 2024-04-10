package userController

import (
	"net/http"

	userService "github.com/WhereNext-co/WhereNext-Backend.git/packages/user/service"
	"github.com/gin-gonic/gin"
)

type UserControllerInterface interface {
	CreateUserInfo(c *gin.Context)
	CheckUserName(c *gin.Context)
	FindUser(c *gin.Context)
	UpdateUserInfo(c *gin.Context)
	IsFriend(c *gin.Context)
	CreateFriendRequest(c *gin.Context)
	AcceptFriendRequest(c *gin.Context)
	DeclineFriendRequest(c *gin.Context)
	CancelFriendRequest(c *gin.Context)
	RemoveFriend(c *gin.Context)
	FriendList(c *gin.Context)
	RequestsSent(c *gin.Context)
	RequestsReceived(c *gin.Context)
}

type UserController struct {
	userService userService.UserServiceInterface
}

func NewUserController(userService userService.UserServiceInterface) *UserController {
	return &UserController{
		userService: userService,
	}
}

// CreateUserInfo at database
func (uc *UserController) CreateUserInfo(c *gin.Context) {
	var user struct {
		UserName       string `json:"userName"`
		Email          string `json:"email"`
		Title          string `json:"title"`
		Name           string `json:"name"`
		Birthdate      string `json:"birthdate"`
		Region         string `json:"region"`
		TelNo          string `json:"telNo"`
		ProfilePicture string `json:"profilePicture"`
		Bio            string `json:"bio"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := uc.userService.CreateUserInfo(user.UserName, user.Email, user.Title, user.Name, user.Birthdate, user.Region, user.TelNo, user.ProfilePicture, user.Bio)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func (uc *UserController) CheckUserName(c *gin.Context) {
	var request struct {
		UserName string `json:"userName"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	exists, err := uc.userService.CheckUserName(request.UserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exists": exists})
}

func (uc *UserController) FindUser(c *gin.Context) {
	var request struct {
		UserName string `json:"userName"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := uc.userService.FindUser(request.UserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (uc *UserController) UpdateUserInfo(c *gin.Context) {
	var user struct {
		UserName       string `json:"userName"`
		Email          string `json:"email"`
		Title          string `json:"title"`
		Name           string `json:"name"`
		Birthdate      string `json:"birthdate"`
		Region         string `json:"region"`
		TelNo          string `json:"telNo"`
		ProfilePicture string `json:"profilePicture"`
		Bio            string `json:"bio"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := uc.userService.UpdateUserInfo(user.UserName, user.Email, user.Title, user.Name, user.Birthdate, user.Region, user.TelNo, user.ProfilePicture, user.Bio)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (uc *UserController) IsFriend(c *gin.Context) {
	var request struct {
		UserName   string `json:"userName"`
		FriendName string `json:"friendName"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	isFriend, err := uc.userService.IsFriend(request.UserName, request.FriendName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"isFriend": isFriend})
}

func (uc *UserController) CreateFriendRequest(c *gin.Context) {
	var request struct {
		UserName   string `json:"userName"`
		FriendName string `json:"friendName"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := uc.userService.CreateFriendRequest(request.UserName, request.FriendName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request created successfully"})
}

func (uc *UserController) AcceptFriendRequest(c *gin.Context) {
	var request struct {
		UserName   string `json:"userName"`
		FriendName string `json:"friendName"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := uc.userService.AcceptFriendRequest(request.UserName, request.FriendName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request accepted successfully"})
}

func (uc *UserController) DeclineFriendRequest(c *gin.Context) {
	var request struct {
		UserName   string `json:"userName"`
		FriendName string `json:"friendName"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := uc.userService.DeclineFriendRequest(request.UserName, request.FriendName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request declined successfully"})
}

func (uc *UserController) CancelFriendRequest(c *gin.Context) {
	var request struct {
		UserName   string `json:"userName"`
		FriendName string `json:"friendName"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := uc.userService.CancelFriendRequest(request.UserName, request.FriendName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request canceled successfully"})
}

func (uc *UserController) RemoveFriend(c *gin.Context) {
	var request struct {
		UserName   string `json:"userName"`
		FriendName string `json:"friendName"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := uc.userService.RemoveFriend(request.UserName, request.FriendName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend removed successfully"})
}

func (uc *UserController) FriendList(c *gin.Context) {
	var request struct {
		UserName string `json:"userName"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	friendList, err := uc.userService.FriendList(request.UserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"friendList": friendList})
}

func (uc *UserController) RequestsSent(c *gin.Context) {
	var request struct {
		UserName string `json:"userName"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	requestsSent, err := uc.userService.RequestsSent(request.UserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"requestsSent": requestsSent})
}

func (uc *UserController) RequestsReceived(c *gin.Context) {
	var request struct {
		UserName string `json:"userName"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	requestsReceived, err := uc.userService.RequestsReceived(request.UserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"requestsReceived": requestsReceived})
}
