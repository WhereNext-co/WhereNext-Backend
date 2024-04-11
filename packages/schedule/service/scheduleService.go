package scheduleService

import (
	"log"
	"time"

	scheduleRepo "github.com/WhereNext-co/WhereNext-Backend.git/packages/schedule/repo"
)

type ScheduleServiceInterface interface {
	CreateLocation(PlaceName string, PlaceGooglePlaceId string, PlaceLocation string,
		PlaceMapLink string, PlacePhotoLink string) error
	CreateSchedule(Name string, Type string,
		Starttime string, Endtime string, Startdate string, Enddate string, Qrcode string, Status string, Location string) (uint, error)
	CreateInvitation(HostID string, InvitedUsers []string) error
	DeleteSchedule(ScheduleID int, HostID int) error
}

type scheduleService struct {
	scheduleRepo scheduleRepo.ScheduleRepoInterface
}

func NewScheduleService(scheduleRepo scheduleRepo.ScheduleRepoInterface) *scheduleService {
	return &scheduleService{scheduleRepo}
}

func (s *scheduleService) CreateLocation(PlaceName string, PlaceGooglePlaceId string, PlaceLocation string, PlaceMapLink string, PlacePhotoLink string) error {
	err := s.scheduleRepo.CreateLocation(PlaceName, PlaceGooglePlaceId, PlaceLocation, PlaceMapLink, PlacePhotoLink)
	if err != nil {
		return err
	}
	return nil
}

func (s *scheduleService) CreateSchedule(Name string, Type string,
	Starttime string, Endtime string, Startdate string, Enddate string, Qrcode string,
	Repeat string, Status string, LocationID string) (uint, error) {
	parsedStartdate, err := time.Parse("2006-01-02", Startdate)
	if err != nil {
		log.Printf("Error parsing starttime: %v", err)
		return 0, err
	}
	parsedEnddate, err := time.Parse("2006-01-02", Enddate)
	if err != nil {
		log.Printf("Error parsing enddate: %v", err)
		return 0, err
	}
	scheduleID, err2 := s.scheduleRepo.CreateSchedule(Name, Type, Starttime, Endtime, parsedStartdate, parsedEnddate, Qrcode, Status, LocationID)
	if err2 != nil {
		log.Printf("Error parsing starttime: %v", err2)
		return 0, err2
	}
	return scheduleID, nil
}

func (s *scheduleService) CreateInvitation(ScheduleID uint, HostID string, InvitedUsers []string) error {
	for _, InviteeID := range InvitedUsers {
		err := s.scheduleRepo.CreateInvitation(ScheduleID, HostID, InviteeID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *scheduleService) DeleteSchedule(ScheduleID int, HostID int) error {
	err := s.scheduleRepo.DeleteSchedule(ScheduleID, HostID)
	if err != nil {
		log.Printf("Error parsing starttime: %v", err)
		return err
	}
	return nil
}
