package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/melfish/br-api/internal/hub"
	"github.com/melfish/br-api/internal/models"
	"github.com/melfish/br-api/internal/service"
)


type NotificationHandler struct {
	svc *service.NotificationService
	hub *hub.Hub
}

func NewNotificationHandler(svc *service.NotificationService, h *hub.Hub) *NotificationHandler {
	return &NotificationHandler{svc: svc, hub: h}
}

// @Summary     Stream notifications via SSE
// @Tags        notifications
// @Produce     text/event-stream
// @Param       recipient_type  query  string  true  "technician or manager"
// @Param       recipient_id    query  string  true  "Recipient UUID"
// @Success     200
// @Failure     400  {object}  ErrorResponse
// @Router      /notifications/stream [get]
func (h *NotificationHandler) Stream(c *gin.Context) {
	recipientType := models.RecipientType(c.Query("recipient_type"))
	if recipientType != models.RecipientTypeTechnician && recipientType != models.RecipientTypeManager {
		Fail(c, http.StatusBadRequest, "recipient_type must be technician or manager")
		return
	}

	recipientID, err := uuid.Parse(c.Query("recipient_id"))
	if err != nil {
		Fail(c, http.StatusBadRequest, "invalid recipient_id")
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	ch := h.hub.Subscribe(recipientType, recipientID)
	defer h.hub.Unsubscribe(recipientType, recipientID, ch)

	for {
		select {
		case <-c.Request.Context().Done():
			return
		case n := <-ch:
			fmt.Fprintf(c.Writer, "data: {\"id\":%q,\"message\":%q,\"type\":%q,\"created_at\":%q}\n\n",
				n.ID, n.Message, n.Type, n.CreatedAt.Format(time.RFC3339))
			c.Writer.Flush()
		}
	}
}

// @Summary     List notifications for a recipient
// @Tags        notifications
// @Produce     json
// @Param       recipient_type  query  string  true  "technician or manager"
// @Param       recipient_id    query  string  true  "Recipient UUID"
// @Success     200  {object}  Response[[]service.NotificationResponse]
// @Failure     400  {object}  ErrorResponse
// @Router      /notifications [get]
func (h *NotificationHandler) List(c *gin.Context) {
	recipientType := models.RecipientType(c.Query("recipient_type"))
	if recipientType != models.RecipientTypeTechnician && recipientType != models.RecipientTypeManager {
		Fail(c, http.StatusBadRequest, "recipient_type must be technician or manager")
		return
	}

	recipientID, err := uuid.Parse(c.Query("recipient_id"))
	if err != nil {
		Fail(c, http.StatusBadRequest, "invalid recipient_id")
		return
	}

	notifications, err := h.svc.List(recipientType, recipientID)
	if err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	Success(c, http.StatusOK, notifications)
}

// @Summary     Mark a notification as read
// @Tags        notifications
// @Param       id            path   string  true  "Notification ID"
// @Param       recipient_id  query  string  true  "Recipient UUID"
// @Success     204
// @Failure     400  {object}  ErrorResponse
// @Router      /notifications/{id}/read [patch]
func (h *NotificationHandler) Read(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		Fail(c, http.StatusBadRequest, "invalid notification id")
		return
	}

	recipientID, err := uuid.Parse(c.Query("recipient_id"))
	if err != nil {
		Fail(c, http.StatusBadRequest, "invalid recipient_id")
		return
	}

	if err := h.svc.Read(id, recipientID); err != nil {
		if errors.Is(err, service.ErrNotificationNotFound) {
			Fail(c, http.StatusNotFound, err.Error())
		} else {
			Fail(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	c.Status(http.StatusNoContent)
}
