package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Kihouww/mini-ai-compute-platform/internal/model"
	"github.com/Kihouww/mini-ai-compute-platform/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func RegisterUserRoutes(r *gin.Engine, userService *service.UserService) {
	h := &UserHandler{userService: userService}

	v1 := r.Group("/v1")
	v1.POST("/users/register", h.Register)
	v1.POST("/users/login", h.Login)
}

func (h *UserHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.userService.Register(c.Request.Context(), &req)
	if errors.Is(err, service.ErrUserAlreadyExists) {
		Error(c, http.StatusConflict, http.StatusConflict, "user already exists")
		return
	}

	if errors.Is(err, service.ErrInvalidInput) {
		Error(c, http.StatusBadRequest, http.StatusBadRequest, err.Error())
		return
	}

	if err != nil {
		Error(c, http.StatusInternalServerError, http.StatusInternalServerError, "internal server error")
		return
	}

	Success(c, resp)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.userService.Login(c.Request.Context(), &req)
	if errors.Is(err, service.ErrInvalidCredentials) {
		Error(c, http.StatusUnauthorized, http.StatusUnauthorized, "invalid username or password")
		return
	}

	if err != nil {
		Error(c, http.StatusInternalServerError, http.StatusInternalServerError, "internal server error")
		return
	}

	Success(c, resp)
}
