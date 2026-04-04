package handler

import (
	"context"
	"time"

	"github.com/highonsemicolon/experiments/appointment-booking/internal/service"
)

type Handler struct {
	coachService   service.CoachService
	bookingService service.BookingService
}

func NewHandler(coachSvc service.CoachService, bookingSvc service.BookingService) *Handler {
	return &Handler{
		coachService:   coachSvc,
		bookingService: bookingSvc,
	}
}


func errPtr(s string) *string { return new(s) }


func (h *Handler) RegisterCoach(ctx context.Context, request RegisterCoachRequestObject) (RegisterCoachResponseObject, error) {
	userID := request.Params.UserID
	
	coach, err := h.coachService.RegisterCoach(userID, request.Body.Name, string(request.Body.Email))
	if err != nil {
		switch err {
		case service.ErrCoachAlreadyExists:
			return RegisterCoach409JSONResponse{ConflictJSONResponse{Error: errPtr(err.Error())}}, nil
		default:
			return RegisterCoach500JSONResponse{InternalErrorJSONResponse{Error: errPtr(err.Error())}}, nil
		}
	}

	return RegisterCoach201JSONResponse(CoachResponse{
		Id:        &coach.ID,
		Name:      &coach.Name,
		Email:     &coach.Email,
		CreatedAt: &coach.CreatedAt,
	}), nil
}


func (h *Handler) SetCoachAvailability(ctx context.Context, request SetCoachAvailabilityRequestObject) (SetCoachAvailabilityResponseObject, error) {
	userID := request.Params.UserID

	tz := ""
	if request.Body.Timezone != nil {
		tz = *request.Body.Timezone
	}

	availability, err := h.coachService.SetAvailability(userID, service.SetAvailabilityInput{
		DayOfWeek: string(request.Body.DayOfWeek),
		StartTime: request.Body.StartTime,
		EndTime:   request.Body.EndTime,
		Timezone:  tz,
	})
	if err != nil {
		switch err {
		case service.ErrCoachNotFound:
			return SetCoachAvailability400JSONResponse{BadRequestJSONResponse{Error: errPtr(err.Error())}}, nil
		default:
			return SetCoachAvailability500JSONResponse{InternalErrorJSONResponse{Error: errPtr(err.Error())}}, nil
		}
	}

	dow := availability.DayOfWeek
	st := availability.StartTime
	et := availability.EndTime
	tz = availability.Timezone
	return SetCoachAvailability201JSONResponse(AvailabilityResponse{
		DayOfWeek: &dow,
		StartTime: &st,
		EndTime:   &et,
		Timezone:  &tz,
	}), nil
}


func (h *Handler) GetCoachAvailability(ctx context.Context, request GetCoachAvailabilityRequestObject) (GetCoachAvailabilityResponseObject, error) {
	availabilities, err := h.coachService.GetCoachAvailability(request.CoachId)
	if err != nil {
		switch err {
		case service.ErrCoachNotFound:
			return GetCoachAvailability404JSONResponse{NotFoundJSONResponse{Error: errPtr(err.Error())}}, nil
		default:
			return GetCoachAvailability500JSONResponse{InternalErrorJSONResponse{Error: errPtr(err.Error())}}, nil
		}
	}

	resp := make(GetCoachAvailability200JSONResponse, 0, len(availabilities))
	for _, a := range availabilities {
		dow := a.DayOfWeek
		st := a.StartTime
		et := a.EndTime
		tz := a.Timezone
		resp = append(resp, AvailabilityResponse{
			DayOfWeek: &dow,
			StartTime: &st,
			EndTime:   &et,
			Timezone:  &tz,
		})
	}
	return resp, nil
}


func (h *Handler) GetAvailableSlots(ctx context.Context, request GetAvailableSlotsRequestObject) (GetAvailableSlotsResponseObject, error) {
	tz := ""
	if request.Params.Timezone != nil {
		tz = *request.Params.Timezone
	}

	date := request.Params.Date.Time.Format("2006-01-02")

	slots, err := h.bookingService.GetAvailableSlots(request.Params.CoachId, date, tz)
	if err != nil {
		switch err {
		case service.ErrCoachNotFound, service.ErrNoAvailability:
			return GetAvailableSlots404JSONResponse{NotFoundJSONResponse{Error: errPtr(err.Error())}}, nil
		case service.ErrInvalidTimezone:
			return GetAvailableSlots400JSONResponse{BadRequestJSONResponse{Error: errPtr(err.Error())}}, nil
		default:
			return GetAvailableSlots500JSONResponse{InternalErrorJSONResponse{Error: errPtr(err.Error())}}, nil
		}
	}

	coachId := request.Params.CoachId
	dateStr := date
	return GetAvailableSlots200JSONResponse(AvailableSlotsResponse{
		CoachId:  &coachId,
		Date:     &dateStr,
		Timezone: &tz,
		Slots:    &slots,
	}), nil
}

func (h *Handler) CreateBooking(ctx context.Context, request CreateBookingRequestObject) (CreateBookingResponseObject, error) {
		userID := request.Params.UserID


	booking, err := h.bookingService.CreateBooking(userID, request.Body.CoachId, request.Body.SlotTime)
	if err != nil {
		switch err {
		case service.ErrCoachNotFound, service.ErrNoAvailability:
			return CreateBooking404JSONResponse{NotFoundJSONResponse{Error: errPtr(err.Error())}}, nil
		case service.ErrSlotOutsideWindow, service.ErrSlotUnavailable:
			return CreateBooking400JSONResponse{BadRequestJSONResponse{Error: errPtr(err.Error())}}, nil
		case service.ErrSlotAlreadyBooked:
			return CreateBooking409JSONResponse{ConflictJSONResponse{Error: errPtr(err.Error())}}, nil
		default:
			return CreateBooking500JSONResponse{InternalErrorJSONResponse{Error: errPtr(err.Error())}}, nil
		}
	}

	id := int(booking.ID)
	status := BookingResponseStatus(booking.Status)
	return CreateBooking201JSONResponse(BookingResponse{
		Id:        &id,
		UserId:    &booking.UserID,
		CoachId:   &booking.CoachID,
		CoachName: &booking.Coach.Name,
		SlotTime:  &booking.StartTime,
		Status:    &status,
		CreatedAt: &booking.CreatedAt,
	}), nil
}

func (h *Handler) GetUserBookings(ctx context.Context, request GetUserBookingsRequestObject) (GetUserBookingsResponseObject, error) {
	userID := request.Params.UserID

	tz := ""
	if request.Params.Timezone != nil {
		tz = *request.Params.Timezone
	}

	bookings, err := h.bookingService.GetMyBookings(userID, tz)
	if err != nil {
		switch err {
		case service.ErrInvalidTimezone:
			return GetUserBookings400JSONResponse{BadRequestJSONResponse{Error: errPtr(err.Error())}}, nil
		default:
			return GetUserBookings500JSONResponse{InternalErrorJSONResponse{Error: errPtr(err.Error())}}, nil
		}
	}

	resp := make(GetUserBookings200JSONResponse, 0, len(bookings))
	for _, b := range bookings {
		id := int(b.ID)
		status := BookingResponseStatus(b.Status)
		slotTime := b.StartTime.UTC().Format(time.RFC3339)
		slotTimeParsed, _ := time.Parse(time.RFC3339, slotTime)
		resp = append(resp, BookingResponse{
			Id:        &id,
			UserId:    &b.UserID,
			CoachId:   &b.CoachID,
			CoachName: &b.Coach.Name,
			SlotTime:  &slotTimeParsed,
			Status:    &status,
			CreatedAt: &b.CreatedAt,
		})
	}
	return resp, nil
}

func (h *Handler) CancelBooking(ctx context.Context, request CancelBookingRequestObject) (CancelBookingResponseObject, error) {
	userID := request.Params.UserID

	if err := h.bookingService.CancelBooking(userID, uint(request.BookingId)); err != nil {
		switch err {
		case service.ErrBookingNotFound:
			return CancelBooking404JSONResponse{NotFoundJSONResponse{Error: errPtr(err.Error())}}, nil
		case service.ErrForbidden:
			return CancelBooking403JSONResponse{ForbiddenJSONResponse{Error: errPtr(err.Error())}}, nil
		default:
			return CancelBooking500JSONResponse{InternalErrorJSONResponse{Error: errPtr(err.Error())}}, nil
		}
	}

	return CancelBooking204Response{}, nil
}