package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/echoed-abyss/qq-like-server/internal/service"
	"github.com/echoed-abyss/qq-like-server/pkg/utils"
)

type MessageHandler struct {
	messageService *service.MessageService
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{
		messageService: service.NewMessageService(),
	}
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	var req service.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	userID := utils.GetUserIDFromContext(c)

	result, err := h.messageService.SendMessage(userID, &req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, result)
}

func (h *MessageHandler) RecallMessage(c *gin.Context) {
	msgIDStr := c.Param("id")
	msgID, err := strconv.ParseUint(msgIDStr, 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "消息ID无效")
		return
	}

	userID := utils.GetUserIDFromContext(c)

	err = h.messageService.RecallMessage(userID, msgID)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *MessageHandler) GetMessageList(c *gin.Context) {
	sessionTypeStr := c.Query("session_type")
	sessionIDStr := c.Query("session_id")
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")

	sessionType, _ := strconv.Atoi(sessionTypeStr)
	sessionID, _ := strconv.ParseUint(sessionIDStr, 10, 64)
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	userID := utils.GetUserIDFromContext(c)
	_ = userID

	messages, total, err := h.messageService.GetMessageList(userID, sessionType, sessionID, page, pageSize)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"list":  messages,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func (h *MessageHandler) AddFavorite(c *gin.Context) {
	var req struct {
		MsgID uint64 `json:"msg_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	userID := utils.GetUserIDFromContext(c)

	err := h.messageService.AddFavorite(userID, req.MsgID)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *MessageHandler) RemoveFavorite(c *gin.Context) {
	favoriteIDStr := c.Param("id")
	favoriteID, err := strconv.ParseUint(favoriteIDStr, 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "收藏ID无效")
		return
	}

	userID := utils.GetUserIDFromContext(c)

	err = h.messageService.RemoveFavorite(userID, favoriteID)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *MessageHandler) GetFavorites(c *gin.Context) {
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	userID := utils.GetUserIDFromContext(c)

	favorites, total, err := h.messageService.GetFavorites(userID, page, pageSize)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"list":  favorites,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}
