package ai

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uniflou/uniflou-lms/backend/internal/shared/pagination"
	"github.com/uniflou/uniflou-lms/backend/internal/shared/response"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) AnalyzeTask(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	attemptID, err := uuid.Parse(c.Param("attempt_id"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	rep, err := h.svc.AnalyzeTask(c.Request.Context(), userID, attemptID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, rep)
}

type trainingSessionRequest struct {
	AttemptIDs []string `json:"attempt_ids" binding:"required,min=1,max=30,dive,uuid"`
	Subject    string   `json:"subject" binding:"required"`
	TaskNumber string   `json:"task_number" binding:"required"`
}

func (h *Handler) AnalyzeTrainingSession(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	var req trainingSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err)
		return
	}
	attemptIDs := make([]uuid.UUID, 0, len(req.AttemptIDs))
	for _, raw := range req.AttemptIDs {
		id, err := uuid.Parse(raw)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid attempt id"})
			return
		}
		attemptIDs = append(attemptIDs, id)
	}
	rep, err := h.svc.AnalyzeTrainingSession(c.Request.Context(), userID, attemptIDs, strings.TrimSpace(req.Subject), strings.TrimSpace(req.TaskNumber))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, rep)
}

func (h *Handler) AnalyzeDashboard(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	rep, err := h.svc.AnalyzeDashboard(c.Request.Context(), userID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, rep)
}

func (h *Handler) AnalyzePlan(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	rep, err := h.svc.AnalyzePlan(c.Request.Context(), userID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, rep)
}

func (h *Handler) AnalyzeProgress(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	rep, err := h.svc.AnalyzeProgress(c.Request.Context(), userID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, rep)
}

func (h *Handler) Reports(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	p := pagination.Parse(c.Query("page"), c.Query("limit"))
	items, total, err := h.svc.Reports(c.Request.Context(), userID, p)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, response.Page{Items: items, Total: total, Page: p.Page})
}
