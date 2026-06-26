package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/echoed-abyss/qq-like-server/internal/service"
	"github.com/echoed-abyss/qq-like-server/pkg/utils"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: service.NewAuthService(),
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Nickname string `json:"nickname"`
}

type LoginRequest struct {
	Account     string `json:"account" binding:"required"`
	Password    string `json:"password" binding:"required"`
	DeviceType  int    `json:"device_type"`
	DeviceName  string `json:"device_name"`
	DeviceModel string `json:"device_model"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	user, err := h.authService.Register(&service.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
	})
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, user)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	clientIP := c.ClientIP()

	result, err := h.authService.Login(&service.LoginRequest{
		Account:  req.Account,
		Password: req.Password,
	}, req.DeviceType, req.DeviceName, req.DeviceModel, clientIP)
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	utils.Success(c, result)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)
	deviceID := c.GetHeader("X-Device-ID")

	err := h.authService.Logout(userID, deviceID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *AuthHandler) Health(c *gin.Context) {
	utils.Success(c, gin.H{
		"status":  "ok",
		"version": "1.0.0",
		"time":    strconv.FormatInt(utils.GetCurrentTimestamp(), 10),
	})
}
