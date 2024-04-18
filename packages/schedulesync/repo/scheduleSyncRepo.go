package scheduleSyncRepo

import (
	"log"
	"time"

	"github.com/WhereNext-co/WhereNext-Backend.git/database" // replace with your actual package name
	"gorm.io/gorm"
)

type ScheduleSyncRepoInterface interface {
    GetFriendsSchedules(uid string, startDate time.Time, endDate time.Time) (map[string][]database.Schedule, error)
	GetSpecificSchedules(uid string, friendUIDs []string, startDate time.Time, endDate time.Time) (map[string][]database.Schedule, error)
}

type scheduleSyncRepo struct {
    dbConn *gorm.DB
}

func NewScheduleSyncRepo(dbConn *gorm.DB) ScheduleSyncRepoInterface {
    return &scheduleSyncRepo{dbConn: dbConn}
}

func (r *scheduleSyncRepo) GetFriendsSchedules(uid string, startDate time.Time, endDate time.Time) (map[string][]database.Schedule, error) {
    var user database.User
    friendSchedules := make(map[string][]database.Schedule)

    // Get the user with the given uid
    err := r.dbConn.First(&user, "uid = ?", uid).Error
    if err != nil {
		log.Printf("Error getting user with uid %s: %v", uid, err)
        return nil, err
    }

    // Get the friends of the user
    var friends []database.User
    err = r.dbConn.Model(&user).Association("Friends").Find(&friends)
    if err != nil {
		log.Printf("Error getting friends of user with uid %s: %v", uid, err)
        return nil, err
    }

    // Add the user to the list of friends to also get the user's own schedules
    friends = append(friends, user)

    // Get the active schedules of the friends (including the user) that are within the given time and date range
    for _, friend := range friends {
        var schedules []database.Schedule
        db := r.dbConn.Joins("LEFT JOIN invitations on invitations.schedule_id = schedules.id").
            Where("(schedules.status = 'Active' and invitations.status = 'Active' and invitations.receiver_uid =?) or (schedules.status = 'Active' and schedules.host_id=?)", friend.Uid, friend.Uid).
            Where("start_date >= ? AND end_date <= ?", startDate, endDate).
            Preload("Invitations", "status='Active'").
            Preload("Invitations.Receiver").
            Find(&schedules)

        if db.Error != nil {
            log.Printf("Error executing schedule query: %v", db.Error)
            return nil, db.Error
        }

        friendSchedules[friend.Uid] = schedules
    }

    return friendSchedules, nil
}

func (r *scheduleSyncRepo) GetSpecificSchedules(uid string, friendUIDs []string, startDate time.Time, endDate time.Time) (map[string][]database.Schedule, error) {
    specificSchedules := make(map[string][]database.Schedule)

	// add user uid to the frienduid
	friendUIDs = append(friendUIDs, uid)
    // Loop through the given friend's UIDs
    for _, friendUID := range friendUIDs {
        var friend database.User

        // Get the friend with the given UID
        err := r.dbConn.First(&friend, "uid = ?", friendUID).Error
        if err != nil {
            log.Printf("Error getting friend with uid %s: %v", friendUID, err)
            return nil, err
        }

        // Get the active schedules of the friend that are within the given time and date range
        var schedules []database.Schedule
        db := r.dbConn.Joins("LEFT JOIN invitations on invitations.schedule_id = schedules.id").
            Where("(schedules.status = 'Active' and invitations.status = 'Active' and invitations.receiver_uid =?) or (schedules.status = 'Active' and schedules.host_id=?)", friend.Uid, friend.Uid).
            Where("start_date >= ? AND end_date <= ?", startDate, endDate).
            Preload("Invitations", "status='Active'").
            Preload("Invitations.Receiver").
            Find(&schedules)

        if db.Error != nil {
            log.Printf("Error executing schedule query: %v", db.Error)
            return nil, db.Error
        }

        // Add the schedules to the map
        specificSchedules[friend.Uid] = schedules
    }

    return specificSchedules, nil
}