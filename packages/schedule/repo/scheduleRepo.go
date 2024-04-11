package scheduleRepo

import (
	"errors"
	"time"

	"github.com/WhereNext-co/WhereNext-Backend.git/database"
	"gorm.io/gorm"
)

type ScheduleRepoInterface interface {
	CreateLocation(PlaceName string, PlaceGooglePlaceId string, PlaceLocation string,
		PlaceMapLink string, PlacePhotoLink string) error
	CreateSchedule(Name string, Type string,
		Starttime string, Endtime string, Startdate time.Time, Enddate time.Time, Qrcode string,
		Status string, Location string) (uint, error)
	CreateInvitation(ScheduleID uint, HostID string, InvitedID string) error
	DeleteSchedule(ScheduleID int, HostID int) error
}

type scheduleRepo struct {
	dbConn *gorm.DB
}

func NewScheduleRepo(dbConn *gorm.DB) *scheduleRepo {
	return &scheduleRepo{dbConn: dbConn}
}

func (s *scheduleRepo) CreateLocation(PlaceName string, PlaceGooglePlaceId string, PlaceLocation string,
	PlaceMapLink string, PlacePhotoLink string) error {
	// Check if location already exists
	var existingLocation database.Location
	if err := s.dbConn.Where("GooglePlaceId = ?", PlaceGooglePlaceId).First(&existingLocation).Error; err == nil {
		return errors.New("location already exists")
	}
	result := s.dbConn.Create(&database.Location{
		Name:          PlaceName,
		GooglePlaceId: PlaceGooglePlaceId,
		Address:       PlaceLocation,
		MapLink:       PlaceMapLink,
		PhotoLink:     PlacePhotoLink,
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *scheduleRepo) CreateSchedule(Name string, Type string,
	Starttime string, Endtime string, Startdate time.Time, Enddate time.Time, Qrcode string,
	Status string, PlaceGooglePlaceId string) (uint, error) {
	var place database.Location
	findlocation := s.dbConn.Where("google_place_id= ?", PlaceGooglePlaceId).First(&place)
	if findlocation.Error != nil {
		return 0, findlocation.Error
	}
	schedule := database.Schedule{Name: Name,
		HostID:     "6",
		StartTime:  Starttime,
		EndTime:    Endtime,
		StartDate:  Startdate,
		EndDate:    Enddate,
		Type:       Type,
		QrCode:     Qrcode,
		LocationID: place.ID}
	result := s.dbConn.Create(&schedule)
	if result.Error != nil {
		return 0, result.Error
	}
	return schedule.ID, nil
}

func (s *scheduleRepo) CreateInvitation(ScheduleID uint, HostID string, InvitedID string) error {
	result := s.dbConn.Create(&database.Invitation{
		ScheduleID:  ScheduleID,
		SenderUid:   HostID,
		ReceiverUid: InvitedID,
		Status:      "pending",
	})
	if result.Error != nil {
		return result.Error
	}
	return nil

}

func (s *scheduleRepo) DeleteSchedule(ScheduleID int, HostID int) error {
	result := s.dbConn.Where("id= ? and host_id= ? ", ScheduleID, HostID).Delete(&database.Schedule{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
