package scheduleRepo

import (
	"log"
	"time"

	"github.com/WhereNext-co/WhereNext-Backend.git/database"
	"gorm.io/gorm"
)

type ScheduleRepoInterface interface {
	CreateLocation(PlaceName string, PlaceGooglePlaceId string, PlaceLocation string,
		PlaceMapLink string, PlacePhotoLink string) error
	FindLocation(PlaceGooglePlaceId string) (database.Location, error)
	FindLocationExist(PlaceGooglePlaceId string) (bool, error)
	FindScheduleByID(scheduleID uint) (database.Schedule, error)
	FindLocationByID(LocationID uint) (database.Location, error)
	ChangeStatusFromDraftToActive(scheduleID uint) error
	CreatePersonalSchedule(HostUid string, Name string, Type string,
		Starttime string, Endtime string, Startdate time.Time, Enddate time.Time, Qrcode string,
		Status string, LocationID uint) error
	CreateRendezvous(HostUid string, Name string, Type string,
		Starttime string, Endtime string, Startdate time.Time, Enddate time.Time, Qrcode string,
		Status string, LocationID uint) (uint, error)
	EditPersonalSchedule(ScheduleID uint, HostUid string, Name string, Type string,
		Starttime string, Endtime string, Startdate time.Time, Enddate time.Time, Qrcode string,
		Status string, LocationID uint) error
	GetActiveSchedule(userID string) ([]database.Schedule, error)
	GetActiveScheduleByDate(userID string, Date time.Time) ([]database.Schedule, error)
	GetDraftRendezvous(userID string) ([]database.Schedule, error)
	GetPastRendezvous(userID string) ([]database.Schedule, error)
	GetPendingRendezvous(userID string) ([]database.Schedule, error)
	GetActiveRendezvous(userID string) ([]database.Schedule, error)
	CreateInvitation(ScheduleID uint, HostID string, InvitedID string) error
	AcceptInvitation(InvitationID uint) error
	RejectInvitation(InvitationID uint) error
	DeleteSchedule(ScheduleID int) error
	RemoveInviteeRendezvous(cheduleID uint, HostID string, InvitedID string) error
}

type scheduleRepo struct {
	dbConn *gorm.DB
}

func NewScheduleRepo(dbConn *gorm.DB) *scheduleRepo {
	return &scheduleRepo{dbConn: dbConn}
}

func (s *scheduleRepo) FindLocation(PlaceGooglePlaceId string) (database.Location, error) {
	var location database.Location
	result := s.dbConn.Where("google_place_id = ?", PlaceGooglePlaceId).First(&location)
	if result.Error != nil {
		return location, result.Error
	}
	return location, nil
}

func (s *scheduleRepo) FindLocationByID(LocationID uint) (database.Location, error) {
	var location database.Location
	result := s.dbConn.Where("id = ?", LocationID).First(&location)
	if result.Error != nil {
		return location, result.Error
	}
	return location, nil
}

func (s *scheduleRepo) FindLocationExist(PlaceGooglePlaceId string) (bool, error) {
	var location database.Location
	result := s.dbConn.Where("google_place_id = ?", PlaceGooglePlaceId).First(&location)
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func (s *scheduleRepo) FindScheduleByID(scheduleID uint) (database.Schedule, error) {
	var schedule database.Schedule
	result := s.dbConn.Where("id = ?", scheduleID).First(&schedule)
	if result.Error != nil {
		return schedule, result.Error
	}
	return schedule, nil
}

func (s *scheduleRepo) CreateLocation(PlaceName string, PlaceGooglePlaceId string, PlaceLocation string,
	PlaceMapLink string, PlacePhotoLink string) error {
	location := database.Location{Name: PlaceName,
		GooglePlaceId: PlaceGooglePlaceId,
		Address:       PlaceLocation,
		MapLink:       PlaceMapLink,
		PhotoLink:     PlacePhotoLink}
	result := s.dbConn.Create(&location)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *scheduleRepo) CreatePersonalSchedule(HostUid string, Name string, Type string,
	Starttime string, Endtime string, Startdate time.Time, Enddate time.Time, Qrcode string,
	Status string, LocationID uint) error {
	schedule := database.Schedule{Name: Name,
		HostID:     HostUid,
		Catagory:   "Scedule",
		StartTime:  Starttime,
		EndTime:    Endtime,
		StartDate:  Startdate,
		EndDate:    Enddate,
		Status:     Status,
		Type:       Type,
		QrCode:     Qrcode,
		LocationID: LocationID}
	result := s.dbConn.Create(&schedule)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *scheduleRepo) EditPersonalSchedule(ScheduleID uint, HostUid string, Name string, Type string,
	Starttime string, Endtime string, Startdate time.Time, Enddate time.Time, Qrcode string,
	Status string, LocationID uint) error {
	var schedule database.Schedule
	schedule.ID = ScheduleID
	schedule.Name = Name
	schedule.HostID = HostUid
	schedule.Catagory = "Schedule"
	schedule.StartTime = Starttime
	schedule.EndTime = Endtime
	schedule.StartDate = Startdate
	schedule.EndDate = Enddate
	schedule.Type = Type
	schedule.Status = Status
	schedule.QrCode = Qrcode
	schedule.LocationID = LocationID
	result := s.dbConn.Save(&schedule)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *scheduleRepo) CreateRendezvous(HostUid string, Name string, Type string,
	Starttime string, Endtime string, Startdate time.Time, Enddate time.Time, Qrcode string,
	Status string, LocationID uint) (uint, error) {
	schedule := database.Schedule{Name: Name,
		HostID:     HostUid,
		Catagory:   "Rendezvous",
		StartTime:  Starttime,
		EndTime:    Endtime,
		StartDate:  Startdate,
		EndDate:    Enddate,
		Status:     Status,
		Type:       Type,
		QrCode:     Qrcode,
		LocationID: LocationID}
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
func (s *scheduleRepo) RemoveInviteeRendezvous(ScheduleID uint, HostID string, InvitedID string) error {
	result := s.dbConn.Where("schedule_id= ? and sender_uid= ? and receiver_uid =?", ScheduleID, HostID, InvitedID).Delete(&database.Invitation{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *scheduleRepo) AcceptInvitation(invitationID uint) error {
	result := s.dbConn.Model(&database.Invitation{}).Where("id = ?", invitationID).Update("status", "Active")
	var invitation database.Invitation
	result = s.dbConn.Where("id=?", invitationID).First(&invitation)
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

func (s *scheduleRepo) ChangeStatusFromDraftToActive(scheduleID uint) error {
	result := s.dbConn.Model(&database.Schedule{}).Where("id = ?", scheduleID).Update("status", "Active")
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

func (s *scheduleRepo) GetActiveSchedule(userID string) ([]database.Schedule, error) {
	var schedules []database.Schedule
	// Query schedules where the user is either the Receiver UID in an active invitation or the host
	//left join เอาฝั่งซ้ายเป็นหลัก คือ schedule
	err := s.dbConn.Joins("LEFT JOIN invitations on invitations.schedule_id = schedules.id").Where("(schedules.status = 'Active' and invitations.status = 'Active' and invitations.receiver_uid =?) or (schedules.status = 'Active' and schedules.host_id=?)", userID, userID).Preload("Location").Preload("Invitations", "status='Active'").Preload("Invitations.Receiver").Find(&schedules)
	if err != nil {
		log.Printf("Error parsing starttime: %v", err)
		return schedules, err.Error
	}
	return schedules, nil
}

func (s *scheduleRepo) GetActiveScheduleByDate(userID string, Date time.Time) ([]database.Schedule, error) {
	var schedules []database.Schedule
	// Query schedules where the user is either the Receiver UID in an active invitation or the host
	//left join เอาฝั่งซ้ายเป็นหลัก คือ schedule

	err := s.dbConn.Joins("LEFT JOIN invitations on invitations.schedule_id = schedules.id").Where("(schedules.status = 'Active' and invitations.status = 'Active' and invitations.receiver_uid =? and schedules.start_date <=? and schedules.end_date >=?) or (schedules.status = 'Active' and schedules.host_id=? and schedules.start_date <=? and schedules.end_date >=?)", userID, Date, Date, userID, Date, Date).Preload("Location").Preload("Invitations", "status='Active'").Preload("Invitations.Receiver").Find(&schedules)
	if err != nil {
		log.Printf("Error parsing starttime: %v", err)
		return schedules, err.Error
	}
	return schedules, nil
}

func (s *scheduleRepo) GetDraftRendezvous(userID string) ([]database.Schedule, error) {
	var schedules []database.Schedule
	err := s.dbConn.Where("schedules.status = 'Draft' and schedules.catagory = 'Rendezvous' and schedules.status = 'Draft' and schedules.host_id=?", userID).Preload("Location").Preload("Invitations").Preload("Invitations.Receiver").Find(&schedules)
	if err != nil {
		log.Printf("Error parsing starttime: %v", err)
		return schedules, err.Error
	}
	return schedules, nil
}

func (s *scheduleRepo) GetPastRendezvous(userID string) ([]database.Schedule, error) {
	var schedules []database.Schedule
	err := s.dbConn.Joins("LEFT JOIN invitations on invitations.schedule_id = schedules.id").Where("(schedules.status = 'Active' and schedules.end_date<now() and schedules.catagory = 'Rendezvous' and invitations.status = 'Active' and invitations.receiver_uid =?) or (schedules.end_date<now() and schedules.catagory = 'Rendezvous' and schedules.status = 'Active' and schedules.host_id=?)", userID, userID).Preload("Location").Preload("Invitations", "status='Active'").Preload("Invitations.Receiver").Find(&schedules)
	if err != nil {
		log.Printf("Error parsing starttime: %v", err)
		return schedules, err.Error
	}
	return schedules, nil
}

func (s *scheduleRepo) GetPendingRendezvous(userID string) ([]database.Schedule, error) {
	var schedules []database.Schedule
	err := s.dbConn.Joins("LEFT JOIN invitations on invitations.schedule_id = schedules.id").Where("schedules.status = 'Active'and schedules.catagory = 'Rendezvous' and invitations.status = 'pending' and invitations.receiver_uid =?", userID).Preload("Location").Preload("Invitations").Preload("Invitations.Receiver").Find(&schedules)
	if err != nil {
		log.Printf("Error parsing starttime: %v", err)
		return schedules, err.Error
	}
	return schedules, nil
}

func (s *scheduleRepo) GetActiveRendezvous(userID string) ([]database.Schedule, error) {
	var schedules []database.Schedule
	err := s.dbConn.Joins("LEFT JOIN invitations on invitations.schedule_id = schedules.id").Where("(schedules.end_date<now() and schedules.catagory='Rendezvous' and schedules.status = 'Active' and invitations.status = 'Active' and invitations.receiver_uid =?) or (schedules.end_date<now() and schedules.catagory='Rendezvous' and schedules.status = 'Active' and schedules.host_id=?)", userID, userID).Preload("Location").Preload("Invitations", "status='Active'").Preload("Invitations.Receiver").Find(&schedules)
	if err != nil {
		log.Printf("Error parsing starttime: %v", err)
		return schedules, err.Error
	}
	return schedules, nil
}
