package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/echoed-abyss/qq-like-server/internal/service"
	"github.com/echoed-abyss/qq-like-server/pkg/utils"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(),
	}
}

func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)

	result, err := h.userService.GetUserInfo(userID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, result)
}

func (h *UserHandler) GetUserProfile(c *gin.Context) {
	targetIDStr := c.Param("id")
	targetID, err := strconv.ParseUint(targetIDStr, 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "用户ID无效")
		return
	}

	currentUserID := utils.GetUserIDFromContext(c)

	result, err := h.userService.GetUserProfile(targetID, currentUserID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, result)
}

func (h *UserHandler) UpdateStatus(c *gin.Context) {
	var req struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	userID := utils.GetUserIDFromContext(c)

	err := h.userService.UpdateOnlineStatus(userID, req.Status)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *UserHandler) GetDevices(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)

	devices, err := h.userService.GetUserDevices(userID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, devices)
}

func (h *UserHandler) KickDevice(c *gin.Context) {
	deviceIDStr := c.Param("id")
	deviceID, err := strconv.ParseUint(deviceIDStr, 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "设备ID无效")
		return
	}

	userID := utils.GetUserIDFromContext(c)

	err = h.userService.KickDevice(userID, deviceID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, nil)
}
