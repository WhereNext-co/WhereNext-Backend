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
	GetScedule(c *gin.Context)
	GetDiary(c *gin.Context)
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
	if schedule.Type == "Rendezvous" {
		err3 := sc.scheduleService.CreateInvitation(scheduleID, schedule.HostID, schedule.InvitedUsers)
		if err3 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create schedule"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Schedule created successfully"})
}

func (sc *ScheduleController) AcceptInvitation(c *gin.Context) {
	var request struct {
		InvitationID uint `json:"InvitationID"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	err := sc.scheduleService.AcceptInvitation(request.InvitationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": " Accepted invitation successfully"})
}

func (sc *ScheduleController) RejectInvitation(c *gin.Context) {
	var request struct {
		InvitationID uint `json:"InvitationID"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	err := sc.scheduleService.RejectInvitation(request.InvitationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": " Reject invitation successfully"})
}

// delete schedule
func (sc *ScheduleController) DeleteSchedule(c *gin.Context) {
	var Deleteschedule struct {
		ScheduleId int `json:"scheduleid"`
	}
	if err := c.ShouldBindJSON(&Deleteschedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	err2 := sc.scheduleService.DeleteSchedule(Deleteschedule.ScheduleId)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Schedule"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Schedule is deleted successfully"})
}

// edit schedule
func (sc *ScheduleController) EditSchedule(c *gin.Context) {
	var schedule struct {
		ScheduleID         uint     `json:"scheduleid"`
		Name               string   `json:"name"`
		Type               string   `json:"type"`
		Starttime          string   `json:"starttime"`
		Endtime            string   `json:"endtime"`
		Startdate          string   `json:"startdate"`
		Enddate            string   `json:"enddate"`
		Qrcode             string   `json:"qrcode"`
		Status             string   `json:"status"`
		User               []string `json:"invitee"`
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
	err2 := sc.scheduleService.EditSchedule(schedule.ScheduleID, schedule.Name, schedule.Type, schedule.Starttime, schedule.Endtime, schedule.Startdate, schedule.Enddate, schedule.Qrcode, schedule.Status, schedule.PlaceName, schedule.PlaceGooglePlaceId, schedule.PlaceLocation, schedule.PlaceMapLink, schedule.PlacePhotoLink)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Update Schedule"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Schedule is update successfully"})
}

func (sc *ScheduleController) GetSchedule(c *gin.Context) {
	var Getschedule struct {
		HostId string `json:"hostid"`
	}
	if err := c.ShouldBindJSON(&Getschedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	scheduleList, err2 := sc.scheduleService.GetSchedule(Getschedule.HostId)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Get Schedule"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Get Schedule successfully", "scheduleList": scheduleList})
}

func (sc *ScheduleController) GetDiary(c *gin.Context) {
	var GetDiary struct {
		HostId string `json:"hostid"`
	}
	if err := c.ShouldBindJSON(&GetDiary); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	diaryList, err2 := sc.scheduleService.GetDiary(GetDiary.HostId)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Get GetDiary"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Get Diary successfully", "diaryList": diaryList})
}
