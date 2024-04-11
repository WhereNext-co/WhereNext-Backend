package scheduleController

import (
	"net/http"

	scheduleService "github.com/WhereNext-co/WhereNext-Backend.git/packages/schedule/service"
	"github.com/gin-gonic/gin"
)

type ScheduleControllerInterface interface {
	CreateSchedule(c *gin.Context)
	EditSchedule(c *gin.Context)
	DeleteSchedule(c *gin.Context)
	CreateInvitation(c *gin.Context)
	AcceptInvitation(c *gin.Context)
	RejectInvitation(c *gin.Context)
}

type ScheduleController struct {
	scheduleService scheduleService.ScheduleServiceInterface
}

func NewScheduleController(scheduleService scheduleService.ScheduleServiceInterface) *ScheduleController {
	return &ScheduleController{
		scheduleService: scheduleService,
	}
}

// CreateSchedule at database
func (sc *ScheduleController) CreateSchedule(c *gin.Context) {
	var schedule struct {
		HostID             string   `json:"hostid"`
		Name               string   `json:"name"`
		Type               string   `json:"type"`
		Starttime          string   `json:"starttime"`
		Endtime            string   `json:"endtime"`
		Startdate          string   `json:"startdate"`
		Enddate            string   `json:"enddate"`
		Qrcode             string   `json:"qrcode"`
		Status             string   `json:"status"`
		InvitedUsers       []string `json:"InvitedUsers"`
		PlaceName          string   `json:"placename"`
		PlaceGooglePlaceId string   `json:"placegoogleplaceid"`
		PlaceLocation      string   `json:"placelocation"`
		PlaceMapLink       string   `json:"placemaplink"`
		PlacePhotoLink     string   `json:"placephotolink"`
	}

	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	err := sc.scheduleService.CreateLocation(schedule.PlaceName, schedule.PlaceGooglePlaceId, schedule.PlaceLocation, schedule.PlaceMapLink, schedule.PlacePhotoLink)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create location"})
		return
	}

	scheduleID, err2 := sc.scheduleService.CreateSchedule(schedule.Name, schedule.Type, schedule.Starttime, schedule.Endtime, schedule.Startdate, schedule.Enddate, schedule.Qrcode, schedule.Status, schedule.PlaceGooglePlaceId)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create schedule"})
		return
	}
	if Type == "Rendezvous" {
		err3 := sc.scheduleService.CreateInvitation(scheduleID, schedule.HostID, schedule.InvitedUsers)
		if err3 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create schedule"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Schedule created successfully"})
}

func (uc *UserController) AcceptInvitation(c *gin.Context) {
	var request struct {
		HostID    string `json:"hostid"`
		InviteeID string `json:"inviteeid"`
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

// delete schedule
func (sc *ScheduleController) DeleteSchedule(c *gin.Context) {
	var Deleteschedule struct {
		ScheduleId int `json:"scheduleid"`
		HostID     int `json:"hostid"`
	}
	if err := c.ShouldBindJSON(&Deleteschedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	err2 := sc.scheduleService.DeleteSchedule(Deleteschedule.ScheduleId, Deleteschedule.HostID)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Schedule"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Schedule is deleted successfully"})
}

// edit schedule
func (sc *ScheduleController) EditSchedule(c *gin.Context) {
	var schedule struct {
		Name               string `json:"name"`
		Type               string `json:"type"`
		Starttime          string `json:"starttime"`
		Endtime            string `json:"endtime"`
		Startdate          string `json:"startdate"`
		Enddate            string `json:"enddate"`
		Qrcode             string `json:"qrcode"`
		Status             string `json:"status"`
		User               array  `json:"status"`
		PlaceName          string `json:"placename"`
		PlaceGooglePlaceId string `json:"placegoogleplaceid"`
		PlaceLocation      string `json:"placelocation"`
		PlaceMapLink       string `json:"placemaplink"`
		PlacePhotoLink     string `json:"placephotolink"`
	}
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	err2 := sc.scheduleService.EditSchedule(schedule.Name, schedule.Type, schedule.Starttime, schedule.Endtime, schedule.Startdate, schedule.Enddate, schedule.Qrcode, schedule.Status, schedule.PlaceName, schedule.PlaceGooglePlaceId, schedule.PlaceLocation, schedule.PlaceMapLink, schedule.PlacePhotoLink)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Edit Schedule"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Schedule is edited successfully"})
}
