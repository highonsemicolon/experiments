package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/highonsemicolon/experiments/appointment-booking/internal/model"
)

type CoachRepository interface {
	RegisterCoach(coach *model.Coach) error
	GetByID(coachID string) (*model.Coach, error)
	CreateAvailability(a *model.Availability) error
	GetAvailabilityByID(coachID string) ([]model.Availability, error)
	GetAvailabilityForDay(coachID string, dayOfWeek string) (*model.Availability, error)
}

var ErrCoachAlreadyExists = errors.New("user is already registered as a coach")

type coachRepository struct {
	db *gorm.DB
}

func NewCoachRepository(db *gorm.DB) CoachRepository {
	return &coachRepository{db: db}
}

func (r *coachRepository) RegisterCoach(coach *model.Coach) error {
	result := r.db.Where(model.Coach{ID: coach.ID}).FirstOrCreate(coach)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrCoachAlreadyExists
	}
	return nil
}

func (r *coachRepository) GetByID(coachID string) (*model.Coach, error) {
	var coach model.Coach
	if err := r.db.First(&coach, "id = ?", coachID).Error; err != nil {
		return nil, err
	}
	return &coach, nil
}

func (r *coachRepository) CreateAvailability(a *model.Availability) error {
    existing := &model.Availability{}
    result := r.db.Where("coach_id = ? AND day_of_week = ?", a.CoachID, a.DayOfWeek).First(existing)
    if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
        return result.Error
    }
    if result.RowsAffected > 0 {
        existing.StartTime = a.StartTime
        existing.EndTime = a.EndTime
        existing.Timezone = a.Timezone
        *a = *existing
        return r.db.Save(a).Error
    }
    return r.db.Create(a).Error
}

func (r *coachRepository) GetAvailabilityByID(coachID string) ([]model.Availability, error) {
	var availabilities []model.Availability
	if err := r.db.Where("coach_id = ?", coachID).Find(&availabilities).Error; err != nil {
		return nil, err
	}
	return availabilities, nil
}

func (r *coachRepository) GetAvailabilityForDay(coachID string, dayOfWeek string) (*model.Availability, error) {
	var availability model.Availability
	if err := r.db.
		Where("coach_id = ? AND day_of_week = ?", coachID, dayOfWeek).
		First(&availability).Error; err != nil {
		return nil, err
	}
	return &availability, nil
}
