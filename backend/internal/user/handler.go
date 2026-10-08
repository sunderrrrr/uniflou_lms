package user

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

func (h *Handler) Me(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	u, profile, err := h.svc.Me(c.Request.Context(), userID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{
		"id":             u.ID,
		"email":          u.Email,
		"role":           u.Role,
		"email_verified": u.EmailVerified,
		"profile":        profile,
	})
}

type patchMeReq struct {
	Grade       *int `json:"grade"`
	TargetScore *int `json:"target_score"`
}

type linkStudentReq struct {
	StudentEmail string `json:"student_email" binding:"required,email"`
}

func (h *Handler) LinkStudent(c *gin.Context) {
	parentID := c.MustGet("userID").(uuid.UUID)
	var req linkStudentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err)
		return
	}
	if err := h.svc.LinkStudent(c.Request.Context(), parentID, req.StudentEmail); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"message": "Ученик привязан"})
}

func (h *Handler) PatchMe(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	var req patchMeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err)
		return
	}
	profile, err := h.svc.UpdateMe(c.Request.Context(), userID, req.Grade, req.TargetScore)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, profile)
}
