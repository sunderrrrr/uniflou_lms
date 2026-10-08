package tutor

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

func (h *Handler) Students(c *gin.Context) {
	items, err := h.svc.Students(c.Request.Context())
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, items)
}

func (h *Handler) StudentDetail(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	data, err := h.svc.StudentDetail(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, data)
}

type commentReq struct {
	Body      string `json:"body" binding:"required"`
	NeedsCall bool   `json:"needs_call"`
}

type targetReq struct {
	University string `json:"university"`
	Specialty  string `json:"specialty"`
	PassScore  int    `json:"pass_score"`
}

func (h *Handler) Comment(c *gin.Context) {
	actorID := c.MustGet("userID").(uuid.UUID)
	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	var req commentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err)
		return
	}
	cmt, err := h.svc.Comment(c.Request.Context(), actorID, studentID, req.Body, req.NeedsCall)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, cmt)
}

func (h *Handler) SetStudentTarget(c *gin.Context) {
	actorID := c.MustGet("userID").(uuid.UUID)
	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	var req targetReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err)
		return
	}
	profile, err := h.svc.SetStudentTarget(c.Request.Context(), actorID, studentID, req.University, req.Specialty, req.PassScore)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, profile)
}

func (h *Handler) RequestHelp(c *gin.Context) {
	actorID := c.MustGet("userID").(uuid.UUID)
	var req struct {
		StudentID uuid.UUID `json:"student_id"`
		Message   string    `json:"message"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err)
		return
	}
	if err := h.svc.RequestHelp(c.Request.Context(), actorID, req.StudentID, req.Message); err != nil {
		response.Fail(c, err)
		return
	}
	c.Status(204)
}

func (h *Handler) Requests(c *gin.Context) {
	tutorID := c.MustGet("userID").(uuid.UUID)
	items, err := h.svc.Requests(c.Request.Context(), tutorID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, items)
}

func (h *Handler) MarkRequestRead(c *gin.Context) {
	tutorID := c.MustGet("userID").(uuid.UUID)
	requestID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	if err := h.svc.MarkRequestRead(c.Request.Context(), tutorID, requestID); err != nil {
		response.Fail(c, err)
		return
	}
	c.Status(204)
}
