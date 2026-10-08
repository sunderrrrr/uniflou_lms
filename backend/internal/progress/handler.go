package progress

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

func (h *Handler) Get(c *gin.Context) {
	studentID := c.MustGet("userID").(uuid.UUID)
	data, err := h.svc.Get(c.Request.Context(), studentID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, data)
}

func (h *Handler) Subjects(c *gin.Context) {
	studentID := c.MustGet("userID").(uuid.UUID)
	items, err := h.svc.Subjects(c.Request.Context(), studentID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, items)
}

func (h *Handler) Topics(c *gin.Context) {
	studentID := c.MustGet("userID").(uuid.UUID)
	items, err := h.svc.Topics(c.Request.Context(), studentID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, items)
}

func (h *Handler) Weak(c *gin.Context) {
	studentID := c.MustGet("userID").(uuid.UUID)
	items, err := h.svc.Weak(c.Request.Context(), studentID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, items)
}

func (h *Handler) MeProgress(c *gin.Context) {
	studentID := c.MustGet("userID").(uuid.UUID)
	dash, err := h.svc.Dashboard(c.Request.Context(), studentID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, dash)
}

type createMockResultReq struct {
	SubjectID  string `json:"subject_id" binding:"required"`
	Title      string `json:"title" binding:"required"`
	TotalTasks int    `json:"total_tasks"`
	Score      int    `json:"score"`
}

func (h *Handler) CreateMockResult(c *gin.Context) {
	studentID := c.MustGet("userID").(uuid.UUID)
	var req createMockResultReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err)
		return
	}
	subjectID, err := uuid.Parse(req.SubjectID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	item, err := h.svc.CreateMockResult(c.Request.Context(), studentID, subjectID, req.Title, req.TotalTasks, req.Score)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, item)
}

func (h *Handler) MockResults(c *gin.Context) {
	studentID := c.MustGet("userID").(uuid.UUID)
	items, err := h.svc.MockResults(c.Request.Context(), studentID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, items)
}
