package model

import "time"

type Coach struct {
    ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Name      string    `gorm:"not null"                json:"name"`
    Email     string    `gorm:"uniqueIndex;not null"    json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

type User struct {
    ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Name      string    `gorm:"not null"                json:"name"`
    Email     string    `gorm:"uniqueIndex;not null"    json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

type Availability struct {
    ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
    CoachID   uint   `gorm:"not null;index"           json:"coach_id"`
    DayOfWeek string `gorm:"not null"                 json:"day_of_week"`
    StartTime string `gorm:"not null"                 json:"start_time"`
    EndTime   string `gorm:"not null"                 json:"end_time"`
    Timezone  string `gorm:"not null;default:'UTC'"   json:"timezone"`

    Coach Coach `gorm:"foreignKey:CoachID" json:"-"`
}

type Booking struct {
    ID        uint      `gorm:"primaryKey;autoIncrement;uniqueIndex:uq_coach_slot" json:"id"`
    CoachID   uint      `gorm:"not null;uniqueIndex:uq_coach_slot"                 json:"coach_id"`
    UserID    uint      `gorm:"not null;index"                                     json:"user_id"`
    StartTime time.Time `gorm:"not null;uniqueIndex:uq_coach_slot"                 json:"start_time"`
    EndTime   time.Time `gorm:"not null"                                           json:"end_time"`
    Status    string    `gorm:"not null;default:'booked'"                          json:"status"`
    CreatedAt time.Time `json:"created_at"`

    Coach Coach `gorm:"foreignKey:CoachID" json:"-"`
    User  User  `gorm:"foreignKey:UserID"  json:"-"`
}
