package scheduleService

import (
	"errors"
	"log"
	"time"

	"github.com/WhereNext-co/WhereNext-Backend.git/database"
	scheduleRepo "github.com/WhereNext-co/WhereNext-Backend.git/packages/schedule/repo"
)

type ScheduleServiceInterface interface {
	CreateLocation(PlaceName string, PlaceGooglePlaceId string, PlaceLocation string,
		PlaceMapLink string, PlacePhotoLink string) error
	CreatePersonalSchedule(HostUid string, Name string, Type string,
		Starttime string, Endtime string, Startdate string, Enddate string, Qrcode string,
		Status string, PlaceName string, PlaceGooglePlaceId string, PlaceLocation string,
		PlaceMapLink string, PlacePhotoLink string) error
	EditPersonalSchedule(ScheduleID uint, HostUid string, Name string, Type string,
		Starttime string, Endtime string, Startdate string, Enddate string, Qrcode string,
		Status string, PlaceName string, PlaceGooglePlaceId string, PlaceLocation string,
		PlaceMapLink string, PlacePhotoLink string) error
	CreateRendezvous(HostUid string, Name string, Type string,
		Starttime string, Endtime string, Startdate string, Enddate string, Qrcode string,
		Status string, PlaceName string, PlaceGooglePlaceId string, PlaceLocation string,
		PlaceMapLink string, PlacePhotoLink string) (uint, error)
	AddInviteeRendezvous(ScheduleID uint, HostID string, InviteeID string) error
	RemoveInviteeRendezvous(ScheduleID uint, HostID string, InviteeID string) error
	GetActiveSchedule(userID string) ([]database.Schedule, error)
	GetActiveScheduleByDate(userID string, Date string) ([]database.Schedule, error)
	GetDraftRendezvous(userID string) ([]database.Schedule, error)
	GetPastRendezvous(userID string) ([]database.Schedule, error)
	GetPendingRendezvous(userID string) ([]database.Schedule, error)
	GetActiveRendezvous(userID string) ([]database.Schedule, error)
	CreateInvitation(ScheduleID uint, HostID string, InvitedUsers []string) error
	AcceptInvitation(InvitationID uint) error
	RejectInvitation(InvitationID uint) error
	DeleteSchedule(ScheduleID int) error
	ChangeStatusFromDraftToActive(scheduleID uint) error
}

type scheduleService struct {
	scheduleRepo scheduleRepo.ScheduleRepoInterface
}

func NewScheduleService(scheduleRepo scheduleRepo.ScheduleRepoInterface) *scheduleService {
	return &scheduleService{scheduleRepo}
}

func (s *scheduleService) CreateLocation(PlaceName string, PlaceGooglePlaceId string, PlaceLocation string, PlaceMapLink string, PlacePhotoLink string) error {
	IsExist, err := s.scheduleRepo.FindLocationExist(PlaceGooglePlaceId)
	if IsExist && err == nil {
		return errors.New("Already Exist")
	}
	err = s.scheduleRepo.CreateLocation(PlaceName, PlaceGooglePlaceId, PlaceLocation, PlaceMapLink, PlacePhotoLink)
	if err != nil {
		return err
	}
	return nil
}

func (s *scheduleService) CreatePersonalSchedule(HostUid string, Name string, Type string,
	Starttime string, Endtime string, Startdate string, Enddate string, Qrcode string,
	Status string, PlaceName string, PlaceGooglePlaceId string, PlaceLocation string,
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
	IsExist, err := s.scheduleRepo.FindLocationExist(PlaceGooglePlaceId)
	if IsExist == false {
		err := s.scheduleRepo.CreateLocation(PlaceName, PlaceGooglePlaceId, PlaceLocation, PlaceMapLink, PlacePhotoLink)
		if err != nil {
			return err
		}
	}
	location, err := s.scheduleRepo.FindLocation(PlaceGooglePlaceId)
	err = s.scheduleRepo.CreatePersonalSchedule(HostUid, Name, Type, Starttime, Endtime, parsedStartdate, parsedEnddate, Qrcode, Status, location.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *scheduleService) CreateRendezvous(HostUid string, Name string, Type string,
	Starttime string, Endtime string, Startdate string, Enddate string, Qrcode string,
	Status string, PlaceName string, PlaceGooglePlaceId string, PlaceLocation string,
	PlaceMapLink string, PlacePhotoLink string) (uint, error) {
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
	IsExist, err := s.scheduleRepo.FindLocationExist(PlaceGooglePlaceId)
	if IsExist == false {
		err := s.scheduleRepo.CreateLocation(PlaceName, PlaceGooglePlaceId, PlaceLocation, PlaceMapLink, PlacePhotoLink)
		if err != nil {
			return 0, err
		}
	}
	location, err := s.scheduleRepo.FindLocation(PlaceGooglePlaceId)
	scheduleID, err := s.scheduleRepo.CreateRendezvous(HostUid, Name, Type, Starttime, Endtime, parsedStartdate, parsedEnddate, Qrcode, Status, location.ID)
	if err != nil {
		return 0, err
	}
	return scheduleID, nil
}

func (s *scheduleService) EditPersonalSchedule(ScheduleID uint, HostUid string, Name string, Type string,
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
	IsExist, err := s.scheduleRepo.FindLocationExist(PlaceGooglePlaceId)
	if IsExist == false {
		err := s.scheduleRepo.CreateLocation(PlaceName, PlaceGooglePlaceId, PlaceLocation, PlaceMapLink, PlacePhotoLink)
		if err != nil {
			return err
		}
		return nil
	}
	location, err := s.scheduleRepo.FindLocation(PlaceGooglePlaceId)
	err = s.scheduleRepo.EditPersonalSchedule(ScheduleID, HostUid, Name, Type, Starttime, Endtime, parsedStartdate, parsedEnddate, Qrcode, Status, location.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *scheduleService) AddInviteeRendezvous(ScheduleID uint, HostID string, InvitedID string) error {
	err := s.scheduleRepo.CreateInvitation(ScheduleID, HostID, InvitedID)
	if err != nil {
		return err
	}
	return nil
}

func (s *scheduleService) RemoveInviteeRendezvous(ScheduleID uint, HostID string, InvitedID string) error {
	err := s.scheduleRepo.RemoveInviteeRendezvous(ScheduleID, HostID, InvitedID)
	if err != nil {
		return err
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

func (s *scheduleService) ChangeStatusFromDraftToActive(ScheduleID uint) error {
	err := s.scheduleRepo.ChangeStatusFromDraftToActive(ScheduleID)
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

func (s *scheduleService) GetActiveSchedule(HostID string) ([]database.Schedule, error) {
	ScheduleList, err := s.scheduleRepo.GetActiveSchedule(HostID)
	if err != nil {
		log.Printf("Error parsing starttime: %v", err)
		return ScheduleList, err
	}
	return ScheduleList, nil
}

func (s *scheduleService) GetActiveScheduleByDate(HostID string, Date string) ([]database.Schedule, error) {
	var schedule []database.Schedule
	parseddate, err := time.Parse("2006-01-02", Date)
	if err != nil {
		log.Printf("Error parsing starttime: %v", err)
		return schedule, err
	}
	ScheduleList, err := s.scheduleRepo.GetActiveScheduleByDate(HostID, parseddate)
	if err != nil {
		log.Printf("Error parsing starttime: %v", err)
		return ScheduleList, err
	}
	return ScheduleList, nil
}

func (s *scheduleService) GetDraftRendezvous(HostID string) ([]database.Schedule, error) {
	ScheduleList, err := s.scheduleRepo.GetDraftRendezvous(HostID)
	if err != nil {
		log.Printf("Error parsing starttime: %v", err)
		return ScheduleList, err
	}
	return ScheduleList, nil
}

func (s *scheduleService) GetPastRendezvous(HostID string) ([]database.Schedule, error) {
	ScheduleList, err := s.scheduleRepo.GetPastRendezvous(HostID)
	if err != nil {
		log.Printf("Error parsing starttime: %v", err)
		return ScheduleList, err
	}
	return ScheduleList, nil
}

func (s *scheduleService) GetPendingRendezvous(HostID string) ([]database.Schedule, error) {
	ScheduleList, err := s.scheduleRepo.GetPendingRendezvous(HostID)
	if err != nil {
		log.Printf("Error parsing starttime: %v", err)
		return ScheduleList, err
	}
	return ScheduleList, nil
}

func (s *scheduleService) GetActiveRendezvous(HostID string) ([]database.Schedule, error) {
	ScheduleList, err := s.scheduleRepo.GetActiveRendezvous(HostID)
	if err != nil {
		log.Printf("Error parsing starttime: %v", err)
		return ScheduleList, err
	}
	return ScheduleList, nil
}
