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
	UserName         string         `gorm:"unique;not null"`
	Email            string         `gorm:"unique"`
	Title            string
	Name             string
	Birthdate        time.Time
	Region           string
	TelNo            string `gorm:"unique;not null"`
	ProfilePicture   string
	Bio              string
	Friends          []User          `gorm:"many2many:UserProfile"`
	RequestsSent     []FriendRequest `gorm:"foreignKey:SenderUid;references:Uid"`
	RequestsReceived []FriendRequest `gorm:"foreignKey:ReceiverUid;references:Uid"`
	RequestSchedules []Invitation    `gorm:"foreignKey:SenderUid;references:Uid"`
	ReceiveSchedules []Invitation    `gorm:"foreignKey:ReceiverUid;references:Uid"`
}

type Schedule struct {
	gorm.Model
	HostID      string
	Host        User `gorm:"foreignKey:HostID"`
	Name        string
	StartTime   string
	EndTime     string
	StartDate   time.Time
	EndDate     time.Time
	Type        string
	QrCode      string
	LocationID  uint
	Invitations []Invitation
}

type Location struct {
	gorm.Model
	Name          string
	GooglePlaceId string
	Address       string
	MapLink       string
	PhotoLink     string
	Schedules     []Schedule
}

type FriendRequest struct {
	gorm.Model
	SenderUid   string `gorm:"type:varchar(255)"`
	Sender      User   `gorm:"foreignKey:SenderUid"`
	ReceiverUid string `gorm:"type:varchar(255)"`
	Receiver    User   `gorm:"foreignKey:ReceiverUid"`
	Status      string
}

type Invitation struct {
	gorm.Model
	ScheduleID  uint     `gorm:"type:varchar(255)"`
	Schedule    Schedule `gorm:"foreignKey:ScheduleID"`
	SenderUid   string   `gorm:"type:varchar(255)"`
	User        User     `gorm:"foreignKey:SenderUid"`
	ReceiverUid string   `gorm:"type:varchar(255)"`
	Receiver    User     `gorm:"foreignKey:ReceiverUid"`
	Status      string
}
