package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/melfish/br-api/internal/service"
)

type ManagerHandler struct {
	svc *service.ManagerService
}

func NewManagerHandler(svc *service.ManagerService) *ManagerHandler {
	return &ManagerHandler{svc: svc}
}

// @Summary     List managers
// @Tags        managers
// @Produce     json
// @Param       page      query  int  false  "Page number"     default(1)
// @Param       page_size query  int  false  "Items per page"  default(20)
// @Success     200  {object}  Response[[]models.Manager]
// @Router      /managers [get]
func (h *ManagerHandler) List(c *gin.Context) {
	page, pageSize := pagination(c)
	result, err := h.svc.List(page, pageSize)
	if err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	SuccessPaged(c, http.StatusOK, result.Items, result.Total, page, pageSize)
}
