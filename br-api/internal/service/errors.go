package service

import "errors"

var (
	ErrConflict             = errors.New("technician has a conflicting job in that time window")
	ErrQuoteNotUnscheduled  = errors.New("quote is not available for scheduling")
	ErrStartsInPast         = errors.New("starts_at must be in the future")
	ErrUnauthorised         = errors.New("not authorised to complete this job")
	ErrJobNotScheduled      = errors.New("job is not in a scheduled state")
	ErrTechnicianNotFound   = errors.New("technician not found")
	ErrNotificationNotFound = errors.New("notification not found")
)
