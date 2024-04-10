package database

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type UserAuth struct {
	gorm.Model
	Email    string
	Password string
	TelNo    string
}

type User struct {
	gorm.Model
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
	RequestsSent     []FriendRequest `gorm:"foreignKey:SenderID"`
	RequestsReceived []FriendRequest `gorm:"foreignKey:ReceiverID"`
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
	IsNotification bool
	ColorCode      string
	Repeat         string
	Weekly         datatypes.JSON
	Monthly        datatypes.JSON
	LocationID     uint
	ParentID       *uint
	Children       []Schedule `gorm:"foreignKey:ParentID"`
	Users          []User     `gorm:"many2many:Invitee"`
}

type Location struct {
	gorm.Model
	Latitude   string
	Longtitude string
	Name       string
	Schedules  []Schedule
}

type Wishlist struct {
	DateTime   time.Time
	ScheduleId uint     `gorm:"primaryKey"`
	Schedule   Schedule `gorm:"foreignKey:ScheduleId"`
	LocationId uint     `gorm:"primaryKey"`
	Location   Location `gorm:"foreignKey:LocationId"`
}

type FriendRequest struct {
	gorm.Model
	SenderID   uint
	Sender     User `gorm:"foreignKey:SenderID"`
	ReceiverID uint
	Receiver   User `gorm:"foreignKey:ReceiverID"`
	Status     string
}
