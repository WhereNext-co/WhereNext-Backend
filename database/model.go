package database

import (
	"time"

	"gorm.io/gorm"
)


type User struct {
    Uid              string `gorm:"primaryKey;type:varchar(255)"`
    CreatedAt        time.Time
    UpdatedAt        time.Time
    DeletedAt        gorm.DeletedAt `gorm:"index"`
    UserName         string `gorm:"unique;not null"`
    Email            string `gorm:"unique"`
    Title            string
    Name             string
    Birthdate        time.Time
    Region           string
    TelNo            string `gorm:"unique;not null"`
    ProfilePicture   string
    Bio              string
    Friends          []User          `gorm:"many2many:UserProfile"`
    Schedules        []Schedule      `gorm:"many2many:Invitee"`
    RequestsSent     []FriendRequest `gorm:"foreignKey:SenderUid;references:Uid"`
    RequestsReceived []FriendRequest `gorm:"foreignKey:ReceiverUid;references:Uid"`
}

type Schedule struct {
	gorm.Model
	Name           string
	StartTime      time.Time
	EndTime        time.Time
	StartDate      time.Time
	EndDate        time.Time
	Type           string
	IsAllDay       bool
	Tag            string
	Note           string
	QrCode         string
	LocationID     uint
	Users          []User     `gorm:"many2many:Invitee;foreignKey:ID;joinForeignKey:ScheduleID;References:Uid;JoinReferences:Uid"`
}

type Location struct {
	gorm.Model
	Latitude   string
	Longtitude string
	Name       string
	Schedules  []Schedule
}

type FriendRequest struct {
    gorm.Model
    SenderUid   string `gorm:"type:varchar(255)"`
    Sender      User `gorm:"foreignKey:SenderUid"`
    ReceiverUid string `gorm:"type:varchar(255)"`
    Receiver    User `gorm:"foreignKey:ReceiverUid"`
    Status      string
}
