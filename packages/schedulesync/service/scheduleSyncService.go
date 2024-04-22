package scheduleSyncService

import (
	"log"
	"time"

	scheduleSyncRepo "github.com/WhereNext-co/WhereNext-Backend.git/packages/schedulesync/repo"
)

type ScheduleSyncServiceInterface interface {
	GetFriendsSchedules(uid string, startTimeStr string, endTimeStr string) (map[string]string, error)
	GetFreeTimeSlots30min(uid string, friendUIDs []string, startDate time.Time, endDate time.Time, duration time.Duration) ([][]time.Time, [][]time.Time, error)
	GetFreeTimeSlotsDaily(uid string, friendUIDs []string, startDate time.Time, endDate time.Time, duration time.Duration) ([][]time.Time, [][]time.Time, error)
}

type scheduleSyncService struct {
	scheduleSyncRepo scheduleSyncRepo.ScheduleSyncRepoInterface
}

func NewScheduleSyncService(scheduleSyncRepo scheduleSyncRepo.ScheduleSyncRepoInterface) *scheduleSyncService {
	return &scheduleSyncService{scheduleSyncRepo}
}

func (s *scheduleSyncService) GetFriendsSchedules(uid string, startTimeStr string, endTimeStr string) (map[string]string, error) {
    layout := time.RFC3339 // This is an example layout. Adjust it to match your Time format.

    startTime, err := time.Parse(layout, startTimeStr)
    if err != nil {
		log.Printf("Error parsing start time: %v", err)
        return nil, err
    }

    endTime, err := time.Parse(layout, endTimeStr)
    if err != nil {
		log.Printf("Error parsing end time: %v", err)
        return nil, err
    }

    friendSchedules, err := s.scheduleSyncRepo.GetFriendsSchedules(uid, startTime, endTime)
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

func (s *scheduleSyncService) GetFreeTimeSlots30min(uid string, friendUIDs []string, startDate time.Time, endDate time.Time, duration time.Duration) ([][]time.Time, [][]time.Time, error) {
    var timePairs [][]time.Time

    for currTime := startDate; currTime.Before(endDate); currTime = currTime.Add(30 * time.Minute) {
        endTime := currTime.Add(duration)
        if endTime.After(endDate) {
            endTime = endDate
        }
        // Check if the duration of the time slot is at least as long as the specified duration
    	if endTime.Sub(currTime) >= duration {
        	timePairs = append(timePairs, []time.Time{currTime, endTime})
    	}
    }

    specificSchedules, err := s.scheduleSyncRepo.GetSpecificSchedules(uid, friendUIDs, startDate, endDate)
    if err != nil {
        log.Printf("Error getting specific schedules: %v", err)
        return nil, nil, err
    }

    var nonOverlappingTimePairs [][]time.Time
    for _, timePair := range timePairs {
        overlaps := false
        for _, specificSchedule := range specificSchedules {
            if (specificSchedule[0].Before(timePair[1]) && specificSchedule[1].After(timePair[0])) ||
                (specificSchedule[0].Equal(timePair[0]) || specificSchedule[1].Equal(timePair[1])) {
                overlaps = true
                break
            }
        }
        if !overlaps {
            nonOverlappingTimePairs = append(nonOverlappingTimePairs, timePair)
        }
    }

    return nonOverlappingTimePairs, specificSchedules, nil
}

func (s *scheduleSyncService) GetFreeTimeSlotsDaily(uid string, friendUIDs []string, startDate time.Time, endDate time.Time, duration time.Duration) ([][]time.Time, [][]time.Time, error) {
    var timePairs [][]time.Time

    for currDate := startDate; currDate.Before(endDate); currDate = currDate.Add(24 * time.Hour) {
        endTime := currDate.Add(duration)
        if endTime.After(endDate) {
            endTime = endDate
        }
        if endTime.Sub(currDate) >= duration {
			timePairs = append(timePairs, []time.Time{currDate, endTime})
		}
    }

    specificSchedules, err := s.scheduleSyncRepo.GetSpecificSchedules(uid, friendUIDs, startDate, endDate)
    if err != nil {
        log.Printf("Error getting specific schedules: %v", err)
        return nil, nil, err
    }

    var nonOverlappingTimePairs [][]time.Time
    for _, timePair := range timePairs {
        overlaps := false
        for _, specificSchedule := range specificSchedules {
            if (specificSchedule[0].Before(timePair[1]) && specificSchedule[1].After(timePair[0])) ||
                (specificSchedule[0].Equal(timePair[0]) || specificSchedule[1].Equal(timePair[1])) {
                overlaps = true
                break
            }
        }
        if !overlaps {
            nonOverlappingTimePairs = append(nonOverlappingTimePairs, timePair)
        }
    }

    return nonOverlappingTimePairs, specificSchedules, nil
}