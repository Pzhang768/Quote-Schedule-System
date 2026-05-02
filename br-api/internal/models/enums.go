package models

type QuoteStatus string

const (
	QuoteStatusUnscheduled QuoteStatus = "unscheduled"
	QuoteStatusScheduled   QuoteStatus = "scheduled"
)

type JobStatus string

const (
	JobStatusScheduled  JobStatus = "scheduled"
	JobStatusCompleted  JobStatus = "completed"
)

type NotificationType string

const (
	NotificationTypeJobAssigned  NotificationType = "job_assigned"
	NotificationTypeJobUpdated   NotificationType = "job_updated"
	NotificationTypeJobCompleted NotificationType = "job_completed"
)