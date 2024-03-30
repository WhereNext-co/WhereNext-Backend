package database

import (
	"gorm.io/gorm"
	"time"
	"gorm.io/datatypes"
)

type UserAuth struct {
	gorm.Model
	Email    string
	Password string
	TelNo    string
}

type User struct{
	gorm.Model
	Title string
	Name string
	Birthdate time.Time
	Region string
	TelNo  string
	ProfilePicture string   
	Friends []User `gorm:"many2many:UserProfile"`
	Schedules []Schedule `gorm:"many2many:Invitee"`
}

type Schedule struct{
	gorm.Model
	Name string
	StartTime time.Time
	EndTime time.Time
	StartDate time.Time
	EndDate time.Time
	Type string 
	IsAllDay bool
	Tag string 
	Note string
	QrCode string
	IsNotification bool
	ColorCode string
	Repeat string
	Weekly datatypes.JSON
	Monthly  datatypes.JSON
	LocationID uint
	ParentID *uint 
	Children []Schedule `gorm:"foreignKey:ParentID"`
	Users []User `gorm:"many2many:Invitee"`
	
}

type Location struct{
	gorm.Model 
	Latitude string 
	Longtitude string 
	Name string
	Schedules []Schedule
}

type Wishlist struct{
	DateTime time.Time
	ScheduleId uint `gorm:"primaryKey"`
	Schedule Schedule `gorm:"foreignKey:ScheduleId"`
	LocationId uint `gorm:"primaryKey"`
	Location Location `gorm:"foreignKey:LocationId"`
}







