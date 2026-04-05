package model

import "time"

type Coach struct {
    ID        string    `gorm:"primaryKey;type:varchar(191)"           json:"id"`
    Name      string    `gorm:"not null;type:varchar(255)"             json:"name"`
    Email     string    `gorm:"uniqueIndex;not null;type:varchar(255)" json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

type Availability struct {
    ID        uint   `gorm:"primaryKey;autoIncrement"                        json:"id"`
    CoachID   string `gorm:"not null;uniqueIndex:uq_coach_day;type:varchar(191)" json:"coach_id"`
    DayOfWeek string `gorm:"not null;uniqueIndex:uq_coach_day;type:varchar(50)"  json:"day_of_week"`
    StartTime string `gorm:"not null;type:varchar(10)"                        json:"start_time"`
    EndTime   string `gorm:"not null;type:varchar(10)"                        json:"end_time"`
    Timezone  string `gorm:"not null;default:'UTC';type:varchar(100)"         json:"timezone"`

    Coach     Coach  `gorm:"foreignKey:CoachID"                               json:"-"`
}

type Booking struct {
    ID        uint      `gorm:"primaryKey;autoIncrement"                              json:"id"`
    CoachID   string    `gorm:"not null;uniqueIndex:uq_coach_slot;type:varchar(191)"  json:"coach_id"`
    UserID    string    `gorm:"not null;index;type:varchar(191)"                      json:"user_id"`
    StartTime time.Time `gorm:"not null;uniqueIndex:uq_coach_slot"                    json:"start_time"`
    EndTime   time.Time `gorm:"not null"                                               json:"end_time"`
    Status    string    `gorm:"not null;default:'booked';type:varchar(50)"            json:"status"`
    CreatedAt time.Time `json:"created_at"`

    Coach     Coach     `gorm:"foreignKey:CoachID"                                    json:"-"`
}