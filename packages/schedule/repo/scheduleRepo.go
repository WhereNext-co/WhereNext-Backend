package scheduleRepo

import (
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
	EditSchedule(ScheduleID uint, Name string, Type string,
		Starttime string, Endtime string, Startdate time.Time, Enddate time.Time, Qrcode string, Status string, PlaceName string, PlaceGooglePlaceId string, PlaceLocation string,
		PlaceMapLink string, PlacePhotoLink string) error
	GetSchedule(HostID string) ([]database.Schedule, error)
	CreateInvitation(ScheduleID uint, HostID string, InvitedID string) error
	AcceptInvitation(InvitationID uint) error
	RejectInvitation(InvitationID uint) error
	DeleteSchedule(ScheduleID int) error
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
	if err := s.dbConn.Where("google_place_id = ?", PlaceGooglePlaceId).First(&existingLocation).Error; err == nil {
		return nil
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
		HostID:     "2",
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

func (s *scheduleRepo) EditSchedule(ScheduleID uint, Name string, Type string,
	Starttime string, Endtime string, Startdate time.Time, Enddate time.Time, Qrcode string, Status string, PlaceName string, PlaceGooglePlaceId string, PlaceLocation string,
	PlaceMapLink string, PlacePhotoLink string) error {
	var place database.Location
	findlocation := s.dbConn.Where("google_place_id= ?", PlaceGooglePlaceId).First(&place)
	if findlocation.Error != nil {
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
		var place2 database.Location
		findlocation := s.dbConn.Where("google_place_id= ?", PlaceGooglePlaceId).First(&place2)
		if findlocation.Error != nil {
			return findlocation.Error
		}
		var schedule database.Schedule
		schedule.Name = Name
		schedule.Type = Type
		schedule.StartTime = Starttime
		schedule.EndTime = Endtime
		schedule.StartDate = Startdate
		schedule.EndDate = Enddate
		schedule.QrCode = Qrcode
		schedule.LocationID = place2.ID
		result1 := s.dbConn.Save(&schedule)
		if result1.Error != nil {
			return result1.Error
		}
		return nil
	} else {
		var schedule database.Schedule
		schedule.Name = Name
		schedule.Type = Type
		schedule.StartTime = Starttime
		schedule.EndTime = Endtime
		schedule.StartDate = Startdate
		schedule.EndDate = Enddate
		schedule.QrCode = Qrcode
		schedule.LocationID = place.ID
		result2 := s.dbConn.Save(&schedule)
		if result2.Error != nil {
			return result2.Error
		}
		return nil
	}
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

func (s *scheduleRepo) AcceptInvitation(invitationID uint) error {
	result := s.dbConn.Model(&database.Invitation{}).Where("id = ?", invitationID).Update("status", "Active")
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *scheduleRepo) RejectInvitation(invitationID uint) error {
	result := s.dbConn.Model(&database.Invitation{}).Where("id = ?", invitationID).Update("status", "Reject")
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *scheduleRepo) DeleteSchedule(ScheduleID int) error {
	result := s.dbConn.Where("id= ? ", ScheduleID).Delete(&database.Schedule{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *scheduleRepo) GetSchedule(HostID string) ([]database.Schedule, error) {
	var schedules []database.Schedule
	result := s.dbConn.Where("host_id = ?", HostID).Find(&schedules)
	if result.Error != nil {
		return nil, result.Error
	}
	return schedules, nil
}

func (s *scheduleRepo) GetDiary(HostID string) ([]database.Schedule, error) {
	var schedules []database.Schedule
	// Join the Schedule table with the Invitation table based on ScheduleID
	// and filter the schedules where the user is either the sender or receiver
	result := s.dbConn.
		Joins("JOIN invitations ON schedules.id = invitations.schedule_id").
		Preload("Invitations").
		Where("invitations.sender_uid = ? OR invitations.receiver_uid = ?", HostID, HostID).
		Find(&schedules)
	if result.Error != nil {
		return nil, result.Error
	}
	return schedules, nil
}
