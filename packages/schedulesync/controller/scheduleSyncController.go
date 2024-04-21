package scheduleSyncController

import (
	"net/http"
	"time"

	scheduleSyncService "github.com/WhereNext-co/WhereNext-Backend.git/packages/schedulesync/service"
	"github.com/gin-gonic/gin"
)

type ScheduleSyncControllerInterface interface {
    GetFriendsSchedules(c *gin.Context)
	GetFreeTimeSlot(c *gin.Context)
}

type scheduleSyncController struct {
	scheduleSyncService scheduleSyncService.ScheduleSyncServiceInterface
}

func NewScheduleSyncController(scheduleSyncService scheduleSyncService.ScheduleSyncServiceInterface) *scheduleSyncController {
	return &scheduleSyncController{
		scheduleSyncService: scheduleSyncService,
	}
}

func (sc *scheduleSyncController) GetFriendsSchedules(c *gin.Context) {
	// Define a struct to represent the request body
	type GetFriendsSchedulesRequest struct {
    	StartTime string `json:"startTime"`
	    EndTime   string `json:"endTime"`
	}

	// Then in your GetFriendsSchedules function:
	var requestBody GetFriendsSchedulesRequest
	if err := c.BindJSON(&requestBody); err != nil {
    	// Handle error
    	return
	}

    uid, exists := c.Get("uid")
    if !exists {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "UID not found"})
        return
    }

    startTimeStr := requestBody.StartTime
    endTimeStr := requestBody.EndTime

    friendAvailability, err := sc.scheduleSyncService.GetFriendsSchedules(uid.(string), startTimeStr, endTimeStr)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, friendAvailability)
}

func (sc *scheduleSyncController) GetFreeTimeSlot(c *gin.Context) {
    var params struct {
        StartTime  string    `json:"startTime"`
        EndTime    string    `json:"endTime"`
        FriendUIDs []string  `json:"friendUIDs"`
        Duration   int64     `json:"duration"` // Duration in seconds
    }

    if err := c.ShouldBindJSON(&params); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    uid, exists := c.Get("uid")
    if !exists {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "UID not found"})
        return
    }

    startDate, err := time.Parse(time.RFC3339, params.StartTime)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time"})
        return
    }

    endDate, err := time.Parse(time.RFC3339, params.EndTime)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end time"})
        return
    }

    duration := time.Duration(params.Duration) * time.Second

    var nonOverlappingSchedules [][]time.Time
	var specificSchedules [][]time.Time
    if duration >= 24*time.Hour {
        nonOverlappingSchedules, specificSchedules, err = sc.scheduleSyncService.GetFreeTimeSlotsDaily(uid.(string), params.FriendUIDs, startDate, endDate, duration)
    } else {
        nonOverlappingSchedules, specificSchedules, err = sc.scheduleSyncService.GetFreeTimeSlots30min(uid.(string), params.FriendUIDs, startDate, endDate, duration)
    }

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
		"nonOverlappingSchedules": nonOverlappingSchedules,
		"specificSchedules": specificSchedules,
	})
}