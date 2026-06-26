package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/echoed-abyss/qq-like-server/internal/model"
	"github.com/echoed-abyss/qq-like-server/internal/service"
	"github.com/echoed-abyss/qq-like-server/pkg/utils"
)

type QRCodeHandler struct {
	qrService *service.QRCodeService
}

func NewQRCodeHandler() *QRCodeHandler {
	return &QRCodeHandler{
		qrService: service.NewQRCodeService(),
	}
}

func (h *QRCodeHandler) GetUserQRCode(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)

	var user model.User
	model.DB.Where("id = ?", userID).First(&user)
	if user.ID == 0 {
		utils.Error(c, http.StatusNotFound, "用户不存在")
		return
	}

	qrUrl, err := h.qrService.GenerateUserQRCode(user.ID, user.QQNumber)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"qr_code": qrUrl,
		"type":    "user",
		"user_id": user.ID,
	})
}

func (h *QRCodeHandler) GetGroupQRCode(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.ParseUint(groupIDStr, 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "群聊ID无效")
		return
	}

	var group model.Group
	model.DB.Where("id = ?", groupID).First(&group)
	if group.ID == 0 {
		utils.Error(c, http.StatusNotFound, "群聊不存在")
		return
	}

	qrUrl, err := h.qrService.GenerateGroupQRCode(group.ID, group.GroupNumber)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"qr_code":   qrUrl,
		"type":      "group",
		"group_id":  group.ID,
		"group_num": group.GroupNumber,
	})
}

func (h *QRCodeHandler) ParseQRCode(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	content, err := h.qrService.ParseQRCode(req.Content)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, content)
}
