package service

import (
	"errors"

	"gorm.io/gorm"

	"github.com/highonsemicolon/experiments/appointment-booking/internal/model"
	"github.com/highonsemicolon/experiments/appointment-booking/internal/repository"
)

var (
	ErrCoachNotFound      = errors.New("coach not found")
	ErrCoachAlreadyExists = errors.New("user is already registered as a coach")
)

type CoachService interface {
	RegisterCoach(userID string, name, email string) (*model.Coach, error)
	GetCoachByID(coachID string) (*model.Coach, error)
	SetAvailability(coachID string, req SetAvailabilityInput) (*model.Availability, error)
	GetCoachAvailability(coachID string) ([]model.Availability, error)
}

type SetAvailabilityInput struct {
	DayOfWeek string
	StartTime string
	EndTime   string
	Timezone  string
}

type coachService struct {
	coachRepo repository.CoachRepository
}

func NewCoachService(coachRepo repository.CoachRepository) CoachService {
	return &coachService{coachRepo: coachRepo}
}

func (s *coachService) RegisterCoach(userID string, name, email string) (*model.Coach, error) {
	coach := &model.Coach{
		ID:    userID,
		Name:  name,
		Email: email,
	}
	if err := s.coachRepo.RegisterCoach(coach); err != nil {
		if errors.Is(err, repository.ErrCoachAlreadyExists) {
			return nil, ErrCoachAlreadyExists
		}
		return nil, err
	}
	return coach, nil
}

func (s *coachService) GetCoachByID(coachID string) (*model.Coach, error) {
	coach, err := s.coachRepo.GetByID(coachID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCoachNotFound
		}
		return nil, err
	}
	return coach, nil
}

func (s *coachService) SetAvailability(coachID string, req SetAvailabilityInput) (*model.Availability, error) {

	tz := req.Timezone
	if tz == "" {
		tz = "UTC"
	}

	availability := &model.Availability{
		CoachID:   coachID,
		DayOfWeek: req.DayOfWeek,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Timezone:  tz,
	}

	if err := s.coachRepo.CreateAvailability(availability); err != nil {
		return nil, err
	}
	return availability, nil
}

func (s *coachService) GetCoachAvailability(coachID string) ([]model.Availability, error) {
	return s.coachRepo.GetAvailabilityByID(coachID)
}