package scheduleController

import (
	"net/http"

	scheduleService "github.com/WhereNext-co/WhereNext-Backend.git/packages/schedule/service"
	"github.com/gin-gonic/gin"
)

type ScheduleControllerInterface interface {
	CreatePersonalSchedule(c *gin.Context)
	CreateRendezvous(c *gin.Context)
	EditPersonalSchedule(c *gin.Context)
	ChangeStatus(c *gin.Context)
	EditRendezvous(c *gin.Context)
	AddInviteeRendezvous(c *gin.Context)
	RemoveInviteeRendezvous(c *gin.Context)
	DeleteSchedule(c *gin.Context)
	GetActiveSchedule(c *gin.Context)
	GetPastRendezvous(c *gin.Context)
	GetPendingRendezvous(c *gin.Context)
	GetDraftRendezvous(c *gin.Context)
	GetActiveRendezvous(c *gin.Context)
	GetActiveScheduleByTime(c *gin.Context)
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

// CreatePersonalSchedule at database
func (sc *ScheduleController) CreatePersonalSchedule(c *gin.Context) {
	var schedule struct {
		HostUid            string `json:"hostuid"`
		Name               string `json:"name"`
		Type               string `json:"type"`
		Starttime          string `json:"starttime"`
		Endtime            string `json:"endtime"`
		Status             string `json:"status"`
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
	err := sc.scheduleService.CreatePersonalSchedule(schedule.HostUid, schedule.Name, schedule.Type, schedule.Starttime, schedule.Endtime, schedule.Status, schedule.PlaceName, schedule.PlaceGooglePlaceId, schedule.PlaceLocation, schedule.PlaceMapLink, schedule.PlacePhotoLink)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create schedule"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Personal Schedule created successfully"})
}

func (sc *ScheduleController) CreateRendezvous(c *gin.Context) {
	var schedule struct {
		HostUid            string   `json:"hostuid"`
		Name               string   `json:"name"`
		Type               string   `json:"type"`
		Starttime          string   `json:"starttime"`
		Endtime            string   `json:"endtime"`
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
	scheduleID, err := sc.scheduleService.CreateRendezvous(schedule.HostUid, schedule.Name, schedule.Type, schedule.Starttime, schedule.Endtime, schedule.Status, schedule.PlaceName, schedule.PlaceGooglePlaceId, schedule.PlaceLocation, schedule.PlaceMapLink, schedule.PlacePhotoLink)

	err = sc.scheduleService.CreateInvitation(scheduleID, schedule.HostUid, schedule.InvitedUsers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create schedule"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Rendervous created successfully"})
}

func (sc *ScheduleController) AcceptInvitation(c *gin.Context) {
	var request struct {
		ScheduleID uint   `json:"scheduleid"`
		ReceiverID string `json:"receiveruid"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	err := sc.scheduleService.AcceptInvitation(request.ScheduleID, request.ReceiverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": " Accepted invitation successfully"})
}

func (sc *ScheduleController) RejectInvitation(c *gin.Context) {
	var request struct {
		ScheduleID uint   `json:"scheduleid"`
		ReceiverID string `json:"receiveruid"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	err := sc.scheduleService.RejectInvitation(request.ScheduleID, request.ReceiverID)
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

// edit Personal schedule
func (sc *ScheduleController) EditPersonalSchedule(c *gin.Context) {
	var schedule struct {
		ScheduleID         uint   `json:"scheduleid"`
		HostID             string `json:"hostuid"`
		Name               string `json:"name"`
		Type               string `json:"type"`
		Starttime          string `json:"starttime"`
		Endtime            string `json:"endtime"`
		Status             string `json:"status"`
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
	err := sc.scheduleService.EditPersonalSchedule(schedule.ScheduleID, schedule.HostID, schedule.Name, schedule.Type, schedule.Starttime, schedule.Endtime, schedule.Status, schedule.PlaceName, schedule.PlaceGooglePlaceId, schedule.PlaceLocation, schedule.PlaceMapLink, schedule.PlacePhotoLink)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Update Schedule"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Schedule is update successfully"})
}

func (sc *ScheduleController) EditRendezvous(c *gin.Context) {
	var schedule struct {
		HostID             string `json:"hostuid"`
		ScheduleID         uint   `json:"scheduleid"`
		Name               string `json:"name"`
		Type               string `json:"type"`
		Starttime          string `json:"starttime"`
		Endtime            string `json:"endtime"`
		Status             string `json:"status"`
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

	err := sc.scheduleService.EditRendezvous(schedule.ScheduleID, schedule.HostID, schedule.Name, schedule.Type, schedule.Starttime, schedule.Endtime, schedule.Status, schedule.PlaceName, schedule.PlaceGooglePlaceId, schedule.PlaceLocation, schedule.PlaceMapLink, schedule.PlacePhotoLink)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Update Schedule"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Schedule is update successfully"})
}

func (sc *ScheduleController) GetActiveSchedule(c *gin.Context) {
	var Getschedule struct {
		HostId string `json:"useruid"`
	}
	if err := c.ShouldBindJSON(&Getschedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	scheduleList, err2 := sc.scheduleService.GetActiveSchedule(Getschedule.HostId)
	personalschedule, rendezvous, err3 := sc.scheduleService.GetAllScheduleMapper(scheduleList)
	if err2 != nil || err3 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Get Schedule"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Get Schedule successfully", "scheduleList": personalschedule, "Rendezvous": rendezvous})
}

func (sc *ScheduleController) GetActiveScheduleByTime(c *gin.Context) {
	var Getschedule struct {
		HostId    string `json:"useruid"`
		StartTime string `json:"starttime"`
		EndTime   string `json:"endtime"`
	}
	if err := c.ShouldBindJSON(&Getschedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	scheduleList, err2 := sc.scheduleService.GetActiveScheduleByTime(Getschedule.HostId, Getschedule.StartTime, Getschedule.EndTime)
	personalschedule, rendezvous, err3 := sc.scheduleService.GetAllScheduleMapper(scheduleList)
	if err2 != nil || err3 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Get Schedule"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Get Schedule successfully", "scheduleList": personalschedule, "Rendezvous": rendezvous})
}

func (sc *ScheduleController) GetDraftRendezvous(c *gin.Context) {
	var Getschedule struct {
		HostId string `json:"useruid"`
	}
	if err := c.ShouldBindJSON(&Getschedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	scheduleList, err2 := sc.scheduleService.GetDraftRendezvous(Getschedule.HostId)
	RendezvousList, err3 := sc.scheduleService.DraftMapper(scheduleList)
	if err3 != nil || err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Get Draft Rendezvous"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Get Draft Rendezvous successfully", "RendezvousList": RendezvousList})
}

func (sc *ScheduleController) GetPastRendezvous(c *gin.Context) {
	var Getschedule struct {
		HostId string `json:"useruid"`
	}
	if err := c.ShouldBindJSON(&Getschedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	scheduleList, err2 := sc.scheduleService.GetPastRendezvous(Getschedule.HostId)
	RendezvousList, err3 := sc.scheduleService.ActiveMapper(scheduleList)
	if err3 != nil || err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Get Past Rendezvous"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Get Past Rendezvous successfully", "RendezvousList": RendezvousList})
}

func (sc *ScheduleController) GetPendingRendezvous(c *gin.Context) {
	var Getschedule struct {
		HostId string `json:"useruid"`
	}
	if err := c.ShouldBindJSON(&Getschedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	scheduleList, err2 := sc.scheduleService.GetPendingRendezvous(Getschedule.HostId)
	RendezvousList, err3 := sc.scheduleService.ActiveMapper(scheduleList)
	if err3 != nil || err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Get Pending Rendezvous"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Get Pending Schedule successfully", "RendezvousList": RendezvousList})
}

func (sc *ScheduleController) GetActiveRendezvous(c *gin.Context) {
	var Getschedule struct {
		HostId string `json:"useruid"`
	}
	if err := c.ShouldBindJSON(&Getschedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	scheduleList, err2 := sc.scheduleService.GetActiveRendezvous(Getschedule.HostId)
	RendezvousList, err3 := sc.scheduleService.ActiveMapper(scheduleList)
	if err3 != nil || err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Get Active Schedule"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Get Active Rendezvous successfully", "RendezvousList": RendezvousList})
}

func (sc *ScheduleController) AddInviteeRendezvous(c *gin.Context) {
	var invitation struct {
		ScheduleID  uint   `json:"scheduleid"`
		SenderUid   string `json:"hostid"`
		ReceicerUid string `json:"inviteeid"`
	}
	if err := c.ShouldBindJSON(&invitation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	err := sc.scheduleService.AddInviteeRendezvous(invitation.ScheduleID, invitation.SenderUid, invitation.ReceicerUid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Add Rendezvous"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Add User to  Rendezvous successfully"})
}

func (sc *ScheduleController) RemoveInviteeRendezvous(c *gin.Context) {
	var invitation struct {
		ScheduleID  uint   `json:"scheduleid"`
		SenderUid   string `json:"hostid"`
		ReceicerUid string `json:"inviteeid"`
	}
	if err := c.ShouldBindJSON(&invitation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	err := sc.scheduleService.RemoveInviteeRendezvous(invitation.ScheduleID, invitation.SenderUid, invitation.ReceicerUid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Remove Rendezvous"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Remove User from Rendezvous successfully"})
}

func (sc *ScheduleController) ChangeStatus(c *gin.Context) {
	var schedule struct {
		ScheduleID uint   `json:"scheduleid"`
		Status     string `json:"status"`
	}
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	err := sc.scheduleService.ChangeStatus(schedule.ScheduleID, schedule.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change status"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Change Status successfully"})
}

func (sc *ScheduleController) AddInviteeRendezvousByID(c *gin.Context) {
	var invitation struct {
		ScheduleID  uint   `json:"scheduleid"`
		ReceiverUid string `json:"useruid"`
	}
	if err := c.ShouldBindJSON(&invitation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	err := sc.scheduleService.AddInviteeRendezvousByID(invitation.ScheduleID, invitation.ReceiverUid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed Add to Rendezvous"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Add User to Rendezvous successfully"})
}
