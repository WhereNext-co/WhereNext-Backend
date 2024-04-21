package scheduleService

import (
	"errors"
	"log"
	"time"

	"github.com/WhereNext-co/WhereNext-Backend.git/database"
	scheduleRepo "github.com/WhereNext-co/WhereNext-Backend.git/packages/schedule/repo"
	userService "github.com/WhereNext-co/WhereNext-Backend.git/packages/user/service"
)

type ScheduleServiceInterface interface {
	CreateLocation(PlaceName string, PlaceGooglePlaceId string, PlaceLocation string,
		PlaceMapLink string, PlacePhotoLink string) error
	CreatePersonalSchedule(HostUid string, Name string, Type string,
		Starttime string, Endtime string,
		Status string, PlaceName string, PlaceGooglePlaceId string, PlaceLocation string,
		PlaceMapLink string, PlacePhotoLink string) error
	EditPersonalSchedule(ScheduleID uint, HostUid string, Name string, Type string,
		Starttime string, Endtime string,
		Status string, PlaceName string, PlaceGooglePlaceId string, PlaceLocation string,
		PlaceMapLink string, PlacePhotoLink string) error
	CreateRendezvous(HostUid string, Name string, Type string,
		Starttime string, Endtime string,
		Status string, PlaceName string, PlaceGooglePlaceId string, PlaceLocation string,
		PlaceMapLink string, PlacePhotoLink string) (uint, error)
	EditRendezvous(ScheduleID uint, HostUid string, Name string, Type string,
		Starttime string, Endtime string,
		Status string, PlaceName string, PlaceGooglePlaceId string, PlaceLocation string,
		PlaceMapLink string, PlacePhotoLink string) error
	AddInviteeRendezvous(ScheduleID uint, HostID string, InviteeID string) error
	AddInviteeRendezvousByID(ScheduleID uint, InviteeID string) error
	RemoveInviteeRendezvous(ScheduleID uint, HostID string, InviteeID string) error
	GetActiveSchedule(userID string) ([]database.Schedule, error)
	GetActiveScheduleByTime(userID string, Starttime string, Endtime string) ([]database.Schedule, error)
	GetDraftRendezvous(userID string) ([]database.Schedule, error)
	GetPastRendezvous(userID string) ([]database.Schedule, error)
	GetPendingRendezvous(userID string) ([]database.Schedule, error)
	GetActiveRendezvous(userID string) ([]database.Schedule, error)
	CreateInvitation(ScheduleID uint, HostID string, InvitedUsers []string) error
	AcceptInvitation(ScheduleID uint, InviteeID string) error
	RejectInvitation(ScheduleID uint, InviteeID string) error
	DeleteSchedule(ScheduleID int) error
	ChangeStatus(scheduleID uint, status string) error
	ActiveMapper(ScheduleList []database.Schedule) ([]Rendezvous, error)
	DraftMapper(ScheduleList []database.Schedule) ([]Rendezvous, error)
	ScheduleMapper(ScheduleList database.Schedule) (Personalschedule, error)
	RendezvousMapper(ScheduleList database.Schedule) (Rendezvous, error)
	GetAllScheduleMapper(ScheduleList []database.Schedule) ([]Personalschedule, []Rendezvous, error)
}

type scheduleService struct {
	scheduleRepo scheduleRepo.ScheduleRepoInterface
	userService  userService.UserServiceInterface
}

func NewScheduleService(scheduleRepo scheduleRepo.ScheduleRepoInterface, userService userService.UserServiceInterface) *scheduleService {
	return &scheduleService{scheduleRepo, userService}
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
	Starttime string, Endtime string,
	Status string, PlaceName string, PlaceGooglePlaceId string, PlaceLocation string,
	PlaceMapLink string, PlacePhotoLink string) error {
	parsedStartTime, err := time.Parse(time.RFC3339, Starttime)
	if err != nil {
		log.Printf("Error parsing starttime: %v", err)
		return err
	}
	parsedEndTime, err := time.Parse(time.RFC3339, Endtime)
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
	err = s.scheduleRepo.CreatePersonalSchedule(HostUid, Name, Type, parsedStartTime, parsedEndTime, Status, location.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *scheduleService) CreateRendezvous(HostUid string, Name string, Type string,
	Starttime string, Endtime string,
	Status string, PlaceName string, PlaceGooglePlaceId string, PlaceLocation string,
	PlaceMapLink string, PlacePhotoLink string) (uint, error) {
	parsedStartTime, err := time.Parse(time.RFC3339, Starttime)
	if err != nil {
		log.Printf("Error parsing starttime: %v", err)
		return 0, err
	}
	parsedEndTime, err := time.Parse(time.RFC3339, Endtime)
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
	scheduleID, err := s.scheduleRepo.CreateRendezvous(HostUid, Name, Type, parsedStartTime, parsedEndTime, Status, location.ID)
	if err != nil {
		return 0, err
	}
	return scheduleID, nil
}

func (s *scheduleService) EditPersonalSchedule(ScheduleID uint, HostUid string, Name string, Type string,
	Starttime string, Endtime string, Status string, PlaceName string, PlaceGooglePlaceId string, PlaceLocation string,
	PlaceMapLink string, PlacePhotoLink string) error {
	parsedStartTime, err := time.Parse(time.RFC3339, Starttime)
	if err != nil {
		log.Printf("Error parsing starttime: %v", err)
		return err
	}
	parsedEndTime, err := time.Parse(time.RFC3339, Endtime)
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
	err = s.scheduleRepo.EditPersonalSchedule(ScheduleID, HostUid, Name, Type, parsedStartTime, parsedEndTime, Status, location.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *scheduleService) EditRendezvous(ScheduleID uint, HostUid string, Name string, Type string,
    Starttime string, Endtime string, Status string, PlaceName string, PlaceGooglePlaceId string, PlaceLocation string,
    PlaceMapLink string, PlacePhotoLink string) error {
    parsedStartTime, err := time.Parse(time.RFC3339, Starttime)
    if err != nil {
        log.Printf("Error parsing starttime: %v", err)
        return err
    }
    parsedEndTime, err := time.Parse(time.RFC3339, Endtime)
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
    err = s.scheduleRepo.EditRendezvous(ScheduleID, HostUid, Name, Type, parsedStartTime, parsedEndTime, Status, location.ID)
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

func (s *scheduleService) AddInviteeRendezvousByID(ScheduleID uint, InvitedID string) error {
	Invitation, err := s.scheduleRepo.FindInvitationByScheduleID(ScheduleID)
	if err != nil {
		return err
	}
	result := s.scheduleRepo.CreateInvitation(ScheduleID, Invitation.SenderUid, InvitedID)
	if result != nil {
		return result
	}
	invitation, err := s.scheduleRepo.FindInvitationByScheduleIDAndInviteeID(ScheduleID, InvitedID)
	err = s.scheduleRepo.AcceptInvitation(invitation.ID)
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

func (s *scheduleService) AcceptInvitation(ScheduleID uint, Invitee string) error {
	invitation, err := s.scheduleRepo.FindInvitationByScheduleIDAndInviteeID(ScheduleID, Invitee)
	err = s.scheduleRepo.AcceptInvitation(invitation.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *scheduleService) RejectInvitation(ScheduleID uint, Invitee string) error {
	invitation, err := s.scheduleRepo.FindInvitationByScheduleIDAndInviteeID(ScheduleID, Invitee)
	err = s.scheduleRepo.RejectInvitation(invitation.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *scheduleService) ChangeStatus(ScheduleID uint, status string) error {
	err := s.scheduleRepo.ChangeStatus(ScheduleID, status)
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

func (s *scheduleService) GetActiveScheduleByTime(HostID string, Starttime string, Endtime string) ([]database.Schedule, error) {
	var schedule []database.Schedule
	parsedstarttime, err := time.Parse(time.RFC3339, Starttime)
	if err != nil {
		log.Printf("Error parsing starttime: %v", err)
		return schedule, err
	}
	parsedendtime, err := time.Parse(time.RFC3339, Endtime)
	if err != nil {
		log.Printf("Error parsing starttime: %v", err)
		return schedule, err
	}
	ScheduleList, err := s.scheduleRepo.GetActiveScheduleByTime(HostID, parsedstarttime, parsedendtime)
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

func (s *scheduleService) GetAllScheduleMapper(ScheduleList []database.Schedule) ([]Personalschedule, []Rendezvous, error) {
	var RendezvousDto []Rendezvous
	var ScheduleDto []Personalschedule
	for _, schedule := range ScheduleList {
		if schedule.Category == "Schedule" {
			result, err := s.ScheduleMapper(schedule)
			if err != nil {
				return []Personalschedule{}, []Rendezvous{}, err
			}
			ScheduleDto = append(ScheduleDto, result)
		} else {
			result, err := s.RendezvousMapper(schedule)
			if err != nil {
				return []Personalschedule{}, []Rendezvous{}, err
			}
			RendezvousDto = append(RendezvousDto, result)
		}
	}
	return ScheduleDto, RendezvousDto, nil
}

func (s *scheduleService) ActiveMapper(ScheduleList []database.Schedule) ([]Rendezvous, error) {
	var RendezvousDto []Rendezvous
	for _, schedule := range ScheduleList {
		var MemberDto []Member
		user, err := s.userService.FindUserByUid(schedule.HostID)
		if err != nil {
			return []Rendezvous{}, err
		}
		MemberDto = append(MemberDto, Member{user.Uid, user.UserName, user.Name, user.ProfilePicture})
		MemberCount := 1
		if len(schedule.Invitations) != 0 {
			for _, invitation := range schedule.Invitations {
				if invitation.Status == "Active" {
					user, err := s.userService.FindUserByUid(invitation.ReceiverUid)
					if err != nil {
						return []Rendezvous{}, err
					}
					MemberDto = append(MemberDto, Member{user.Uid, user.UserName, user.Name, user.ProfilePicture})
					MemberCount += 1
				}

			}
		}
		RendezvousDto = append(RendezvousDto, Rendezvous{schedule.ID, schedule.HostID, schedule.Name, schedule.Type, schedule.Category, schedule.StartTime, schedule.EndTime, schedule.Status, schedule.Location.Name, schedule.Location.GooglePlaceId, schedule.Location.Address, schedule.Location.MapLink, schedule.Location.PhotoLink, MemberCount, MemberDto})
	}
	return RendezvousDto, nil
}

func (s *scheduleService) DraftMapper(ScheduleList []database.Schedule) ([]Rendezvous, error) {
	var RendezvousDto []Rendezvous
	for _, schedule := range ScheduleList {
		var MemberDto []Member
		user, err := s.userService.FindUserByUid(schedule.HostID)
		if err != nil {
			return []Rendezvous{}, err
		}
		MemberDto = append(MemberDto, Member{user.Uid, user.UserName, user.Name, user.ProfilePicture})
		MemberCount := 1
		if len(schedule.Invitations) != 0 {
			for _, invitation := range schedule.Invitations {
				if invitation.Status == "pending" {
					user, err := s.userService.FindUserByUid(invitation.ReceiverUid)
					if err != nil {
						return []Rendezvous{}, err
					}
					MemberDto = append(MemberDto, Member{user.Uid, user.UserName, user.Name, user.ProfilePicture})
					MemberCount += 1
				}

			}
		}
		RendezvousDto = append(RendezvousDto, Rendezvous{schedule.ID, schedule.HostID, schedule.Name, schedule.Type, schedule.Category, schedule.StartTime, schedule.EndTime, schedule.Status, schedule.Location.Name, schedule.Location.GooglePlaceId, schedule.Location.Address, schedule.Location.MapLink, schedule.Location.PhotoLink, MemberCount, MemberDto})
	}
	return RendezvousDto, nil
}

func (s *scheduleService) ScheduleMapper(schedule database.Schedule) (Personalschedule, error) {
	ScheduleDto := Personalschedule{
		ScheduleID:         schedule.ID,
		HostUid:            schedule.HostID,
		Name:               schedule.Name,
		Type:               schedule.Type,
		Category:           schedule.Category,
		Starttime:          schedule.StartTime,
		Endtime:            schedule.EndTime,
		Status:             schedule.Status,
		PlaceName:          schedule.Location.Name,
		PlaceGooglePlaceId: schedule.Location.GooglePlaceId,
		PlaceLocation:      schedule.Location.Address,
		PlaceMapLink:       schedule.Location.MapLink,
		PlacePhotoLink:     schedule.Location.PhotoLink,
	}
	return ScheduleDto, nil
}

func (s *scheduleService) RendezvousMapper(schedule database.Schedule) (Rendezvous, error) {
	var MemberDto []Member
	user, err := s.userService.FindUserByUid(schedule.HostID)
	if err != nil {
		return Rendezvous{}, err
	}
	MemberDto = append(MemberDto, Member{user.Uid, user.UserName, user.Name, user.ProfilePicture})
	MemberCount := 1
	if len(schedule.Invitations) != 0 {
		for _, invitation := range schedule.Invitations {
			if invitation.Status == "Active" {
				user, err := s.userService.FindUserByUid(invitation.ReceiverUid)
				if err != nil {
					return Rendezvous{}, err
				}
				MemberDto = append(MemberDto, Member{user.Uid, user.UserName, user.Name, user.ProfilePicture})
				MemberCount += 1
			}

		}
	}
	RendezvousDto := Rendezvous{
		ScheduleID:         schedule.ID,
		HostUid:            schedule.HostID,
		Name:               schedule.Name,
		Type:               schedule.Type,
		Category:           schedule.Category,
		Starttime:          schedule.StartTime,
		Endtime:            schedule.EndTime,
		Status:             schedule.Status,
		PlaceName:          schedule.Location.Name,
		PlaceGooglePlaceId: schedule.Location.GooglePlaceId,
		PlaceLocation:      schedule.Location.Address,
		PlaceMapLink:       schedule.Location.MapLink,
		PlacePhotoLink:     schedule.Location.PhotoLink,
		MemberCount:        MemberCount,
		Member:             MemberDto,
	}
	return RendezvousDto, nil
}

type Member struct {
	UserID         string `json:"Useruid"`
	UserName       string `json:"userName"`
	Name           string `json:"name"`
	ProfilePicture string `json:"profilePicture"`
}

type Rendezvous struct {
	ScheduleID         uint      `json:"scheduleid"`
	HostUid            string    `json:"hostuid"`
	Name               string    `json:"name"`
	Type               string    `json:"type"`
	Category           string    `json:"category"`
	Starttime          time.Time `json:"starttime"`
	Endtime            time.Time `json:"endtime"`
	Status             string    `json:"status"`
	PlaceName          string    `json:"placename"`
	PlaceGooglePlaceId string    `json:"placegoogleplaceid"`
	PlaceLocation      string    `json:"placelocation"`
	PlaceMapLink       string    `json:"placemaplink"`
	PlacePhotoLink     string    `json:"placephotolink"`
	MemberCount        int       `json:"membercount"`
	Member             []Member  `json:"member"`
}

type Personalschedule struct {
	ScheduleID         uint      `json:"scheduleid"`
	HostUid            string    `json:"hostuid"`
	Name               string    `json:"name"`
	Type               string    `json:"type"`
	Category           string    `json:"category"`
	Starttime          time.Time `json:"starttime"`
	Endtime            time.Time `json:"endtime"`
	Status             string    `json:"status"`
	PlaceName          string    `json:"placename"`
	PlaceGooglePlaceId string    `json:"placegoogleplaceid"`
	PlaceLocation      string    `json:"placelocation"`
	PlaceMapLink       string    `json:"placemaplink"`
	PlacePhotoLink     string    `json:"placephotolink"`
}
