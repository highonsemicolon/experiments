package repository

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/highonsemicolon/experiments/appointment-booking/internal/model"
)

type BookingRepository interface {
	BeginTx() *gorm.DB
	GetBookedSlotsForDay(coachID uint, from, to time.Time) ([]time.Time, error)
	LockAndGetSlot(tx *gorm.DB, coachID uint, startTime time.Time) (*model.Booking, error)
	CreateBookingTx(tx *gorm.DB, booking *model.Booking) error
	GetBookingByID(bookingID uint) (*model.Booking, error)
	GetUserBookings(userID uint) ([]model.Booking, error)
	CancelBooking(bookingID, userID uint) error
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) BeginTx() *gorm.DB {
	return r.db.Begin()
}

func (r *bookingRepository) GetBookedSlotsForDay(coachID uint, from, to time.Time) ([]time.Time, error) {
	var slots []time.Time
	err := r.db.Model(&model.Booking{}).
		Where("coach_id = ? AND start_time >= ? AND start_time < ? AND status = 'booked'",
			coachID, from, to).
		Pluck("start_time", &slots).Error
	return slots, err
}

func (r *bookingRepository) LockAndGetSlot(tx *gorm.DB, coachID uint, startTime time.Time) (*model.Booking, error) {
	var booking model.Booking
	err := tx.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("coach_id = ? AND start_time = ? AND status = 'booked'", coachID, startTime).
		First(&booking).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepository) CreateBookingTx(tx *gorm.DB, booking *model.Booking) error {
	return tx.Create(booking).Error
}

func (r *bookingRepository) GetBookingByID(bookingID uint) (*model.Booking, error) {
	var booking model.Booking
	if err := r.db.Preload("Coach").First(&booking, bookingID).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepository) GetUserBookings(userID uint) ([]model.Booking, error) {
	var bookings []model.Booking
	if err := r.db.Preload("Coach").
		Where("user_id = ? AND status = 'booked'", userID).
		Order("start_time ASC").
		Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *bookingRepository) CancelBooking(bookingID, userID uint) error {
	result := r.db.Model(&model.Booking{}).
		Where("id = ? AND user_id = ?", bookingID, userID).
		Update("status", "cancelled")

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
