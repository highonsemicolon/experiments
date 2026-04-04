package service

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/highonsemicolon/experiments/appointment-booking/internal/model"
	"github.com/highonsemicolon/experiments/appointment-booking/internal/repository"
)

var (
	ErrSlotUnavailable    = errors.New("slot is not available")
	ErrSlotAlreadyBooked  = errors.New("slot is already booked")
	ErrBookingNotFound    = errors.New("booking not found")
	ErrForbidden          = errors.New("you do not own this booking")
	ErrInvalidTimezone    = errors.New("invalid timezone")
	ErrNoAvailability     = errors.New("coach has no availability on this day")
	ErrSlotOutsideWindow  = errors.New("slot is outside coach availability window")
)

type BookingService interface {
	GetAvailableSlots(coachID string, date string, timezone string) ([]time.Time, error)
	CreateBooking(userID string, coachID string, slotTime time.Time) (*model.Booking, error)
	GetMyBookings(userID string, timezone string) ([]model.Booking, error)
	CancelBooking(userID string, bookingID uint) error
}

type bookingService struct {
	bookingRepo repository.BookingRepository
	coachRepo   repository.CoachRepository
}

func NewBookingService(
	bookingRepo repository.BookingRepository,
	coachRepo repository.CoachRepository,
) BookingService {
	return &bookingService{
		bookingRepo: bookingRepo,
		coachRepo:   coachRepo,
	}
}


func (s *bookingService) GetAvailableSlots(coachID string, date string, timezone string) ([]time.Time, error) {
	loc, err := resolveTimezone(timezone)
	if err != nil {
		return nil, ErrInvalidTimezone
	}
	requestedDate, err := time.ParseInLocation("2006-01-02", date, loc)
	if err != nil {
		return nil, fmt.Errorf("invalid date format, use YYYY-MM-DD")
	}

	dayOfWeek := requestedDate.Weekday().String() 
	
	availability, err := s.coachRepo.GetAvailabilityForDay(coachID, dayOfWeek)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoAvailability
		}
		return nil, err
	}
	coachLoc, err := resolveTimezone(availability.Timezone)
	if err != nil {
		return nil, fmt.Errorf("coach has invalid timezone configured")
	}

	windowStart, windowEnd, err := buildWindow(requestedDate, availability.StartTime, availability.EndTime, coachLoc)
	if err != nil {
		return nil, err
	}
	bookedSlots, err := s.bookingRepo.GetBookedSlotsForDay(coachID, windowStart.UTC(), windowEnd.UTC())
	if err != nil {
		return nil, err
	}

	bookedSet := make(map[time.Time]struct{}, len(bookedSlots))
	for _, t := range bookedSlots {
		bookedSet[t.UTC()] = struct{}{}
	}

	var slots []time.Time
	for slot := windowStart; slot.Before(windowEnd); slot = slot.Add(30 * time.Minute) {
		slotUTC := slot.UTC()
		if _, booked := bookedSet[slotUTC]; !booked {
			slots = append(slots, slotUTC.In(loc))
		}
	}

	return slots, nil
}

func (s *bookingService) CreateBooking(userID string, coachID string, slotTime time.Time) (*model.Booking, error) {
	slotUTC := slotTime.UTC()

	coach, err := s.coachRepo.GetByID(coachID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCoachNotFound
		}
		return nil, err
	}

	dayOfWeek := slotUTC.Weekday().String()
	availability, err := s.coachRepo.GetAvailabilityForDay(coachID, dayOfWeek)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoAvailability
		}
		return nil, err
	}

	coachLoc, _ := resolveTimezone(availability.Timezone)
	windowStart, windowEnd, err := buildWindow(slotUTC.In(coachLoc), availability.StartTime, availability.EndTime, coachLoc)
	if err != nil {
		return nil, err
	}
	if slotUTC.Before(windowStart.UTC()) || slotUTC.Equal(windowEnd.UTC()) || slotUTC.After(windowEnd.UTC()) {
		return nil, ErrSlotOutsideWindow
	}

	if slotUTC.Minute()%30 != 0 || slotUTC.Second() != 0 {
		return nil, ErrSlotUnavailable
	}

	tx := s.bookingRepo.BeginTx()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	existing, err := s.bookingRepo.LockAndGetSlot(tx, coachID, slotUTC)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if existing != nil {
		tx.Rollback()
		return nil, ErrSlotAlreadyBooked
	}

	booking := &model.Booking{
		UserID:    userID,
		CoachID:   coachID,
		StartTime: slotUTC,
		EndTime:   slotUTC.Add(30 * time.Minute),
		Status:    "booked",
	}

	if err := s.bookingRepo.CreateBookingTx(tx, booking); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	booking.Coach = *coach
	return booking, nil
}

func (s *bookingService) GetMyBookings(userID string, timezone string) ([]model.Booking, error) {
	loc, err := resolveTimezone(timezone)
	if err != nil {
		return nil, ErrInvalidTimezone
	}

	bookings, err := s.bookingRepo.GetUserBookings(userID)
	if err != nil {
		return nil, err
	}

	for i := range bookings {
		bookings[i].StartTime = bookings[i].StartTime.In(loc)
		bookings[i].EndTime = bookings[i].EndTime.In(loc)
	}

	return bookings, nil
}

func (s *bookingService) CancelBooking(userID string, bookingID uint) error {
	booking, err := s.bookingRepo.GetBookingByID(bookingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrBookingNotFound
		}
		return err
	}

	if booking.UserID != userID {
		return ErrForbidden
	}

	return s.bookingRepo.CancelBooking(bookingID, userID)
}
