package scheduleSyncService

import (
	"log"
	"time"

	scheduleSyncRepo "github.com/WhereNext-co/WhereNext-Backend.git/packages/schedulesync/repo"
)

type ScheduleSyncServiceInterface interface {
	GetFriendsSchedules(uid string, startDateStr string, endDateStr string) (map[string]string, error)
}

type scheduleSyncService struct {
	scheduleSyncRepo scheduleSyncRepo.ScheduleSyncRepoInterface
}

func NewScheduleSyncService(scheduleSyncRepo scheduleSyncRepo.ScheduleSyncRepoInterface) *scheduleSyncService {
	return &scheduleSyncService{scheduleSyncRepo}
}

func (s *scheduleSyncService) GetFriendsSchedules(uid string, startDateStr string, endDateStr string) (map[string]string, error) {
    layout := "2006-01-02T15:04:05Z" // This is an example layout. Adjust it to match your date format.

    startDate, err := time.Parse(layout, startDateStr)
    if err != nil {
		log.Printf("Error parsing start date: %v", err)
        return nil, err
    }

    endDate, err := time.Parse(layout, endDateStr)
    if err != nil {
		log.Printf("Error parsing end date: %v", err)
        return nil, err
    }

    friendSchedules, err := s.scheduleSyncRepo.GetFriendsSchedules(uid, startDate, endDate)
    if err != nil {
		log.Printf("Error getting friends' schedules: %v", err)
        return nil, err
    }

    friendAvailability := make(map[string]string)
    for uid, schedules := range friendSchedules {
        if len(schedules) == 0 {
            friendAvailability[uid] = "available"
        } else {
            friendAvailability[uid] = "busy"
        }
    }

    return friendAvailability, nil
}