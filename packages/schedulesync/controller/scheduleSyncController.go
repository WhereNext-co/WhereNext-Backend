package scheduleSyncController

import (
	"net/http"

	scheduleSyncService "github.com/WhereNext-co/WhereNext-Backend.git/packages/schedulesync/service"
	"github.com/gin-gonic/gin"
)

type ScheduleSyncControllerInterface interface {
    GetFriendsSchedules(c *gin.Context)
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
    uid := c.Query("uid")
    startDateStr := c.Query("startDate")
    endDateStr := c.Query("endDate")

    friendAvailability, err := sc.scheduleSyncService.GetFriendsSchedules(uid, startDateStr, endDateStr)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, friendAvailability)
}