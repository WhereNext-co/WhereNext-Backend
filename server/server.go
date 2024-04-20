package server

import (
	"log"
	"os"

	"github.com/WhereNext-co/WhereNext-Backend.git/database"
	"github.com/WhereNext-co/WhereNext-Backend.git/middleware"
	auth "github.com/WhereNext-co/WhereNext-Backend.git/packages/auth"
	authController "github.com/WhereNext-co/WhereNext-Backend.git/packages/auth/controller"
	authService "github.com/WhereNext-co/WhereNext-Backend.git/packages/auth/service"
	userController "github.com/WhereNext-co/WhereNext-Backend.git/packages/user/controller"
	userRepo "github.com/WhereNext-co/WhereNext-Backend.git/packages/user/repo"
	userService "github.com/WhereNext-co/WhereNext-Backend.git/packages/user/service"

	scheduleController "github.com/WhereNext-co/WhereNext-Backend.git/packages/schedule/controller"
	scheduleRepo "github.com/WhereNext-co/WhereNext-Backend.git/packages/schedule/repo"
	scheduleService "github.com/WhereNext-co/WhereNext-Backend.git/packages/schedule/service"
	scheduleSyncController "github.com/WhereNext-co/WhereNext-Backend.git/packages/schedulesync/controller"
	scheduleSyncRepo "github.com/WhereNext-co/WhereNext-Backend.git/packages/schedulesync/repo"
	scheduleSyncService "github.com/WhereNext-co/WhereNext-Backend.git/packages/schedulesync/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/twilio/twilio-go"
)

func InitServer() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Fail to .env file")
	}
	dbConn := database.InitDB()
	userRepo := userRepo.NewUserRepo(dbConn)
	scheduleRepo := scheduleRepo.NewScheduleRepo(dbConn)
	scheduleSyncRepo := scheduleSyncRepo.NewScheduleSyncRepo(dbConn)
	// Initialize Firebase
	authClient, err := auth.InitializeFirebase()
	if err != nil {
		log.Fatalf("error initializing Firebase: %v", err)
	}

	// Initialize Twilio
	twilioAccountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	twilioAuthToken := os.Getenv("TWILIO_AUTH_TOKEN")
	twilioClient := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: twilioAccountSid,
		Password: twilioAuthToken,
	})

	// Initialize the Authentication services and controllers
	authService := authService.NewAuthService(authClient, twilioClient)
	authController := authController.NewAuthController(authService)
	// Initialize the User services and controllers
	userService := userService.NewUserService(userRepo)
	userController := userController.NewUserController(userService)
	// Initialize the Schedule services and controllers
	scheduleService := scheduleService.NewScheduleService(scheduleRepo, userService)
	scheduleController := scheduleController.NewScheduleController(scheduleService)
	// Initialize the Schedule Sync services and controllers
	scheduleSyncService := scheduleSyncService.NewScheduleSyncService(scheduleSyncRepo)
	scheduleSyncController := scheduleSyncController.NewScheduleSyncController(scheduleSyncService)
	r := gin.Default()
	r.Use(cors.Default())

	// Auth routes
	r.POST("/auth/createFirebaseUser", authController.CreateFirebaseUser)
	r.POST("/auth/updateFirebaseUserPassword", authController.UpdateFirebaseUserPassword)

	// User routes
	r.POST("/users/create-info", middleware.VerifyToken(), userController.CreateUserInfo)
	r.POST("/users/check-username", userController.CheckUserName)
	r.POST("/users/get-profile", userController.FindUserByUid)
	r.PUT("/users/profile", userController.UpdateUserInfo)
	// Friend routes
	r.POST("/users/get-friends", userController.FriendList)
	r.POST("/users/friends/isfriend", userController.IsFriend)
	r.DELETE("/users/friends", userController.RemoveFriend)
	r.POST("/users/friends/friendinfo", userController.FindFriendInfo)
	// Friend request routes
	r.POST("/users/friendrequest", userController.CreateFriendRequest)
	r.PUT("/users/friendrequest", userController.AcceptFriendRequest)
	r.DELETE("/users/friendrequest/decline", userController.DeclineFriendRequest)
	r.DELETE("/users/friendrequest/cancel", userController.CancelFriendRequest)
	r.POST("/users/get-friendrequest", userController.RequestsReceived)

	// Schedule routes
	r.POST("/schedules/create-personalschedule", scheduleController.CreatePersonalSchedule)
	r.DELETE("/schedules/delete-schedule", scheduleController.DeleteSchedule)
	r.PUT("/schedules/edit-schedule", scheduleController.EditPersonalSchedule)
	r.PATCH("/schedules/change-status", scheduleController.ChangeStatus)
	r.POST("/schedules/get-allschedule", scheduleController.GetActiveSchedule)
	r.POST("/schedules/get-schedulebytime", scheduleController.GetActiveScheduleByTime)
	// Rendezvous routes
	r.POST("/rendezvous/create-rendezvous", scheduleController.CreateRendezvous)
	r.POST("/rendezvous/get-draft-rendezvous", scheduleController.GetDraftRendezvous)
	r.POST("/rendezvous/get-past-rendezvous", scheduleController.GetPastRendezvous)
	r.POST("/rendezvous/get-active-rendezvous", scheduleController.GetActiveRendezvous)
	r.POST("/rendezvous/get-pending-rendezvous", scheduleController.GetPendingRendezvous)
	r.PUT("/rendezvous/edit-rendezvous", scheduleController.EditRendezvous)
	r.POST("/rendezvous/add-user-rendezvous", scheduleController.AddInviteeRendezvous)
	r.DELETE("/rendezvous/remove-user-rendezvous", scheduleController.RemoveInviteeRendezvous)
	r.POST("/rendezvous/add-user-byscheduleID", scheduleController.AddInviteeRendezvousByID)
	// Invitation routes
	r.PATCH("/rendezvous/accept-invitation", scheduleController.AcceptInvitation)
	r.PATCH("/rendezvous/reject-invitation", scheduleController.RejectInvitation)
	// Schedule Sync routes
	r.POST("/schedulesync/get-friends-schedules", scheduleSyncController.GetFriendsSchedules)
	r.POST("/schedulesync/get-free-timeslot", scheduleSyncController.GetFreeTimeSlot)
	port := os.Getenv("PORT")
	r.Run(":" + port)
}
