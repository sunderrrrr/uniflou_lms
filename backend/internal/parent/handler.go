package parent

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uniflou/uniflou-lms/backend/internal/shared/response"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Children(c *gin.Context) {
	parentID := c.MustGet("userID").(uuid.UUID)
	items, err := h.svc.Children(c.Request.Context(), parentID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, items)
}

func (h *Handler) ChildProgress(c *gin.Context) {
	parentID := c.MustGet("userID").(uuid.UUID)
	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	data, err := h.svc.ChildProgress(c.Request.Context(), parentID, studentID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) ChildReports(c *gin.Context) {
	parentID := c.MustGet("userID").(uuid.UUID)
	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	data, err := h.svc.ChildReports(c.Request.Context(), parentID, studentID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, data)
}
