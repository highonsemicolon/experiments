package repository

import (
	"gorm.io/gorm"

	"github.com/highonsemicolon/experiments/appointment-booking/internal/model"
)

type CoachRepository interface {
	GetByID(coachID uint) (*model.Coach, error)
	CreateAvailability(a *model.Availability) error
	GetAvailabilityByID(coachID uint) ([]model.Availability, error)
	GetAvailabilityForDay(coachID uint, dayOfWeek string) (*model.Availability, error)
}

type coachRepository struct {
	db *gorm.DB
}

func NewCoachRepository(db *gorm.DB) CoachRepository {
	return &coachRepository{db: db}
}

func (r *coachRepository) GetByID(coachID uint) (*model.Coach, error) {
	var coach model.Coach
	if err := r.db.First(&coach, coachID).Error; err != nil {
		return nil, err
	}
	return &coach, nil
}

func (r *coachRepository) CreateAvailability(a *model.Availability) error {
	return r.db.
		Where(model.Availability{CoachID: a.CoachID, DayOfWeek: a.DayOfWeek}).
		Assign(model.Availability{
			StartTime: a.StartTime,
			EndTime:   a.EndTime,
			Timezone:  a.Timezone,
		}).
		FirstOrCreate(a).Error
}

func (r *coachRepository) GetAvailabilityByID(coachID uint) ([]model.Availability, error) {
	var availabilities []model.Availability
	if err := r.db.Where("coach_id = ?", coachID).Find(&availabilities).Error; err != nil {
		return nil, err
	}
	return availabilities, nil
}

func (r *coachRepository) GetAvailabilityForDay(coachID uint, dayOfWeek string) (*model.Availability, error) {
	var availability model.Availability
	if err := r.db.
		Where("coach_id = ? AND day_of_week = ?", coachID, dayOfWeek).
		First(&availability).Error; err != nil {
		return nil, err
	}
	return &availability, nil
}
