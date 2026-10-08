package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uniflou/uniflou-lms/backend/internal/shared/pagination"
	"github.com/uniflou/uniflou-lms/backend/internal/shared/response"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Users(c *gin.Context) {
	p := pagination.Parse(c.Query("page"), c.Query("limit"))
	items, total, err := h.svc.Users(c.Request.Context(), p)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, response.Page{Items: items, Total: total, Page: p.Page})
}

func (h *Handler) UpdateUserRole(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	var req struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err)
		return
	}
	actorID := c.MustGet("userID").(uuid.UUID)
	if err := h.svc.UpdateUserRole(c.Request.Context(), actorID, userID, req.Role); err != nil {
		response.Fail(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	actorID := c.MustGet("userID").(uuid.UUID)
	if err := h.svc.DeleteUser(c.Request.Context(), actorID, userID); err != nil {
		response.Fail(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) Subjects(c *gin.Context) {
	items, err := h.svc.Subjects(c.Request.Context())
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, items)
}

func (h *Handler) CreateSubject(c *gin.Context) {
	var input SubjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, err)
		return
	}
	item, err := h.svc.CreateSubject(c.Request.Context(), input)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, item)
}

func (h *Handler) CreateTopic(c *gin.Context) {
	var input TopicInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, err)
		return
	}
	item, err := h.svc.CreateTopic(c.Request.Context(), input)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, item)
}

func (h *Handler) Tasks(c *gin.Context) {
	p := pagination.Parse(c.Query("page"), c.Query("limit"))
	items, total, err := h.svc.Tasks(c.Request.Context(), p)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, response.Page{Items: items, Total: total, Page: p.Page})
}

func (h *Handler) CreateTask(c *gin.Context) {
	var input TaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, err)
		return
	}
	item, err := h.svc.CreateTask(c.Request.Context(), input)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, item)
}

func (h *Handler) UpdateTask(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	var input TaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, err)
		return
	}
	item, err := h.svc.UpdateTask(c.Request.Context(), id, input)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, item)
}

func (h *Handler) DeleteTask(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	if err := h.svc.DeleteTask(c.Request.Context(), id); err != nil {
		response.Fail(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func parsePageLimit(value string, fallback int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 {
		return fallback
	}
	return parsed
}
