package scheduleService

import (
	"log"
	"time"

	"github.com/WhereNext-co/WhereNext-Backend.git/database"
	scheduleRepo "github.com/WhereNext-co/WhereNext-Backend.git/packages/schedule/repo"
)

type ScheduleServiceInterface interface {
	CreateLocation(PlaceName string, PlaceGooglePlaceId string, PlaceLocation string,
		PlaceMapLink string, PlacePhotoLink string) error
	CreateSchedule(Name string, Type string,
		Starttime string, Endtime string, Startdate string, Enddate string, Qrcode string, Status string, Location string) (uint, error)
	EditSchedule(ScheduleID uint, Name string, Type string,
		Starttime string, Endtime string, Startdate string, Enddate string, Qrcode string, Status string, PlaceName string, PlaceGooglePlaceId string, PlaceLocation string,
		PlaceMapLink string, PlacePhotoLink string) error
	GetSchedule(HostID string) ([]database.Schedule, error)
	GetDiary(HostID string) ([]database.Schedule, error)
	CreateInvitation(ScheduleID uint, HostID string, InvitedUsers []string) error
	AcceptInvitation(InvitationID uint) error
	RejectInvitation(InvitationID uint) error
	DeleteSchedule(ScheduleID int) error
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
	Status string, LocationID string) (uint, error) {
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

func (s *scheduleService) EditSchedule(ScheduleID uint, Name string, Type string,
	Starttime string, Endtime string, Startdate string, Enddate string, Qrcode string, Status string, PlaceName string, PlaceGooglePlaceId string, PlaceLocation string,
	PlaceMapLink string, PlacePhotoLink string) error {
	parsedStartdate, err := time.Parse("2006-01-02", Startdate)
	if err != nil {
		log.Printf("Error parsing starttime: %v", err)
		return err
	}
	parsedEnddate, err := time.Parse("2006-01-02", Enddate)
	if err != nil {
		log.Printf("Error parsing enddate: %v", err)
		return err
	}
	err2 := s.scheduleRepo.EditSchedule(ScheduleID, Name, Type, Starttime, Endtime, parsedStartdate, parsedEnddate, Qrcode, Status, PlaceName, PlaceGooglePlaceId, PlaceLocation, PlaceMapLink, PlacePhotoLink)
	if err2 != nil {
		return err2
	}
	return nil
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

func (s *scheduleService) AcceptInvitation(InvitationID uint) error {
	err := s.scheduleRepo.AcceptInvitation(InvitationID)
	if err != nil {
		return err
	}
	return nil
}

func (s *scheduleService) RejectInvitation(InvitationID uint) error {
	err := s.scheduleRepo.RejectInvitation(InvitationID)
	if err != nil {
		return err
	}
	return nil
}

func (s *scheduleService) DeleteSchedule(ScheduleID int) error {
	err := s.scheduleRepo.DeleteSchedule(ScheduleID)
	if err != nil {
		log.Printf("Error parsing starttime: %v", err)
		return err
	}
	return nil
}

func (s *scheduleService) GetSchedule(HostID string) ([]database.Schedule, error) {
	ScheduleList, err := s.scheduleRepo.GetSchedule(HostID)
	if err != nil {
		log.Printf("Error parsing starttime: %v", err)
		return ScheduleList, err
	}
	return ScheduleList, nil
}

func (s *scheduleService) GetDiary(HostID string) ([]database.Schedule, error) {
	DiaryList, err := s.scheduleRepo.GetSchedule(HostID)
	if err != nil {
		log.Printf("Error parsing starttime: %v", err)
		return DiaryList, err
	}
	return DiaryList, nil
}
