package task

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	apperrors "github.com/uniflou/uniflou-lms/backend/internal/shared/errors"
	"github.com/uniflou/uniflou-lms/backend/internal/shared/pagination"
	"github.com/uniflou/uniflou-lms/backend/internal/shared/response"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) List(c *gin.Context) {
	p := pagination.Parse(c.Query("page"), c.Query("limit"))
	var topicID *uuid.UUID
	if tid := c.Query("topic_id"); tid != "" {
		id, err := uuid.Parse(tid)
		if err == nil {
			topicID = &id
		}
	}
	items, total, err := h.svc.List(c.Request.Context(), topicID, p)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, response.Page{Items: items, Total: total, Page: p.Page})
}

func (h *Handler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	t, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, t)
}

func (h *Handler) Next(c *gin.Context) {
	studentID := c.MustGet("userID").(uuid.UUID)
	limit := 1
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}
	t, err := h.svc.Next(c.Request.Context(), studentID, c.Query("mode"), c.Query("subject_id"), c.Query("topic_id"), limit)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, t)
}

func (h *Handler) NextBatch(c *gin.Context) {
	studentID := c.MustGet("userID").(uuid.UUID)
	limit := 5
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}
	items, err := h.svc.NextBatch(c.Request.Context(), studentID, c.Query("mode"), c.Query("subject_id"), c.Query("topic_id"), limit)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, items)
}

func (h *Handler) TrainerCatalog(c *gin.Context) {
	items, err := h.svc.TrainerCatalog(c.Request.Context())
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, items)
}

func (h *Handler) TrainerTopics(c *gin.Context) {
	items, err := h.svc.TrainerTopics(c.Request.Context())
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, items)
}

func (h *Handler) MockSession(c *gin.Context) {
	subjectID, err := uuid.Parse(c.Query("subject_id"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	items, err := h.svc.MockSession(c.Request.Context(), subjectID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, items)
}

func (h *Handler) TrainingBatch(c *gin.Context) {
	subjectID, err := uuid.Parse(c.Query("subject_id"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	taskNumber, err := strconv.Atoi(c.Query("task_number"))
	if err != nil || taskNumber < 1 {
		response.Fail(c, fmt.Errorf("%w: номер задания должен быть положительным числом", apperrors.ErrValidation))
		return
	}
	limit := 3
	if value := c.Query("limit"); value != "" {
		if parsed, parseErr := strconv.Atoi(value); parseErr == nil && parsed > 0 {
			limit = parsed
		}
	}
	items, err := h.svc.TrainingBatch(c.Request.Context(), subjectID, taskNumber, limit)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, items)
}

type attemptReq struct {
	Answer    string `json:"answer" binding:"required"`
	TimeSpent int    `json:"time_spent"`
}

func (h *Handler) Attempt(c *gin.Context) {
	studentID := c.MustGet("userID").(uuid.UUID)
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	var req attemptReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err)
		return
	}
	res, err := h.svc.SubmitAttempt(c.Request.Context(), studentID, taskID, req.Answer, req.TimeSpent)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) History(c *gin.Context) {
	studentID := c.MustGet("userID").(uuid.UUID)
	p := pagination.Parse(c.Query("page"), c.Query("limit"))
	items, total, err := h.svc.History(c.Request.Context(), studentID, p)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, response.Page{Items: items, Total: total, Page: p.Page})
}
