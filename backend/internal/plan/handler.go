package plan

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	apperrors "github.com/uniflou/uniflou-lms/backend/internal/shared/errors"
	"github.com/uniflou/uniflou-lms/backend/internal/shared/response"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetPlan(c *gin.Context) {
	studentID := c.MustGet("userID").(uuid.UUID)
	data, err := h.svc.GetPlan(c.Request.Context(), studentID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) PutPlan(c *gin.Context) {
	studentID := c.MustGet("userID").(uuid.UUID)
	var req struct {
		Universities []PlanSelection `json:"universities"`
		Deadlines    []struct {
			Title string `json:"title"`
			Date  string `json:"date"`
			Type  string `json:"type"`
		} `json:"deadlines"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err)
		return
	}
	deadlines := make([]DeadlineInput, 0, len(req.Deadlines))
	for _, deadline := range req.Deadlines {
		parsedDate, err := time.Parse("2006-01-02", deadline.Date)
		if err != nil {
			response.Fail(c, fmt.Errorf("%w: неверный формат даты", apperrors.ErrValidation))
			return
		}
		deadlines = append(deadlines, DeadlineInput{Title: deadline.Title, Date: parsedDate, Type: deadline.Type})
	}
	if err := h.svc.SavePlan(c.Request.Context(), studentID, req.Universities, deadlines); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"message": "План сохранён"})
}

func (h *Handler) ListUniversities(c *gin.Context) {
	items, err := h.svc.Universities(c.Request.Context())
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, items)
}

type addUniReq struct {
	UniversityID uuid.UUID `json:"university_id" binding:"required"`
	Priority     string    `json:"priority"`
}

func (h *Handler) AddUniversity(c *gin.Context) {
	studentID := c.MustGet("userID").(uuid.UUID)
	var req addUniReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err)
		return
	}
	sp, err := h.svc.AddUniversity(c.Request.Context(), studentID, req.UniversityID, req.Priority)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, sp)
}

func (h *Handler) DeleteUniversity(c *gin.Context) {
	studentID := c.MustGet("userID").(uuid.UUID)
	planID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	if err := h.svc.RemoveUniversity(c.Request.Context(), studentID, planID); err != nil {
		response.Fail(c, err)
		return
	}
	c.Status(204)
}

func (h *Handler) Deadlines(c *gin.Context) {
	studentID := c.MustGet("userID").(uuid.UUID)
	items, err := h.svc.Deadlines(c.Request.Context(), studentID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, items)
}

func (h *Handler) AddDeadline(c *gin.Context) {
	studentID := c.MustGet("userID").(uuid.UUID)
	var req struct {
		Title string `json:"title" binding:"required"`
		Date  string `json:"date" binding:"required"`
		Type  string `json:"type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err)
		return
	}
	item, err := h.svc.AddDeadline(c.Request.Context(), studentID, req.Title, req.Date, req.Type)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, item)
}

func (h *Handler) DeleteDeadline(c *gin.Context) {
	studentID := c.MustGet("userID").(uuid.UUID)
	deadlineID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	if err := h.svc.RemoveDeadline(c.Request.Context(), studentID, deadlineID); err != nil {
		response.Fail(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
