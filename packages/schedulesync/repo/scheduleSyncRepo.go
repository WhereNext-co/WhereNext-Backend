package scheduleSyncRepo

import (
	"log"
	"time"

	"github.com/WhereNext-co/WhereNext-Backend.git/database" // replace with your actual package name
	"gorm.io/gorm"
)

type ScheduleSyncRepoInterface interface {
    GetFriendsSchedules(uid string, startTime time.Time, endTime time.Time) (map[string]UserSchedules, error)
	GetSpecificSchedules(uid string, friendUIDs []string, startDate time.Time, endDate time.Time) ([][]time.Time, error)
}

type scheduleSyncRepo struct {
    dbConn *gorm.DB
}

func NewScheduleSyncRepo(dbConn *gorm.DB) ScheduleSyncRepoInterface {
    return &scheduleSyncRepo{dbConn: dbConn}
}

type UserSchedules struct {
    User      database.User
    Schedules []database.Schedule
}

func (r *scheduleSyncRepo) GetFriendsSchedules(uid string, startTime time.Time, endTime time.Time) (map[string]UserSchedules, error) {
    var user database.User
    friendSchedules := make(map[string]UserSchedules)

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
            Where("start_time <= ? AND end_time >= ?", endTime, startTime).
            Preload("Invitations", "status='Active'").
            Preload("Invitations.Receiver").
            Distinct().
            Find(&schedules)

        if db.Error != nil {
            log.Printf("Error executing schedule query: %v", db.Error)
            return nil, db.Error
        }

        friendSchedules[friend.Uid] = UserSchedules{
            User:      friend,
            Schedules: schedules,
        }
    }

    return friendSchedules, nil
}

func (r *scheduleSyncRepo) GetSpecificSchedules(uid string, friendUIDs []string, startDate time.Time, endDate time.Time) ([][]time.Time, error) {
    var specificSchedules [][]time.Time

    friendUIDs = append(friendUIDs, uid)
    for _, friendUID := range friendUIDs {
        var friend database.User

        err := r.dbConn.First(&friend, "uid = ?", friendUID).Error
        if err != nil {
            log.Printf("Error getting friend with uid %s: %v", friendUID, err)
            return nil, err
        }

        var schedules []database.Schedule
        db := r.dbConn.Joins("LEFT JOIN invitations on invitations.schedule_id = schedules.id").
            Where("(schedules.status = 'Active' and invitations.status = 'Active' and invitations.receiver_uid =?) or (schedules.status = 'Active' and schedules.host_id=?)", friend.Uid, friend.Uid).
            Where("start_time <= ? AND end_time >= ?", endDate, startDate).
            Preload("Invitations", "status='Active'").
            Preload("Invitations.Receiver").
			Distinct().
            Find(&schedules)

        if db.Error != nil {
            log.Printf("Error executing schedule query: %v", db.Error)
            return nil, db.Error
        }

        for _, schedule := range schedules {
            specificSchedules = append(specificSchedules, []time.Time{schedule.StartTime, schedule.EndTime})
        }
    }

    return specificSchedules, nil
}