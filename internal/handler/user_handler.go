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

func (h *UserHandler) GetFriendList(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)

	result, err := h.userService.GetFriendGroups(userID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, result)
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

func (h *UserHandler) CheckIn(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)

	expGain, level, err := h.userService.CheckIn(userID)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"exp_gain": expGain,
		"level":    level,
	})
}

func (h *UserHandler) LikeUser(c *gin.Context) {
	targetIDStr := c.Param("id")
	targetID, err := strconv.ParseUint(targetIDStr, 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "用户ID无效")
		return
	}

	likes, err := h.userService.LikeUser(targetID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, gin.H{"likes": likes})
}

func (h *UserHandler) GetLikeRank(c *gin.Context) {
	users, err := h.userService.GetLikeRank()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, users)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req struct {
		Nickname  string `json:"nickname"`
		Signature string `json:"signature"`
		Bio       string `json:"bio"`
		Tags      string `json:"tags"`
		Gender    int    `json:"gender"`
		Age       int    `json:"age"`
		Avatar    string `json:"avatar"`
		Banner    string `json:"banner"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	userID := utils.GetUserIDFromContext(c)

	err := h.userService.UpdateProfile(userID, req.Nickname, req.Signature, req.Bio, req.Tags, req.Gender, req.Age, req.Avatar, req.Banner)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *UserHandler) UpdatePassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6,max=32"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	userID := utils.GetUserIDFromContext(c)

	err := h.userService.UpdatePassword(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *UserHandler) DeleteAccount(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)

	err := h.userService.DeleteAccount(userID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, nil)
}
