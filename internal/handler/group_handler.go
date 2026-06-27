package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/echoed-abyss/qq-like-server/internal/service"
	"github.com/echoed-abyss/qq-like-server/pkg/utils"
)

type GroupHandler struct {
	groupService *service.GroupService
}

func NewGroupHandler() *GroupHandler {
	return &GroupHandler{
		groupService: service.NewGroupService(),
	}
}

func (h *GroupHandler) GetGroupInfo(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.ParseUint(groupIDStr, 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "群聊ID无效")
		return
	}

	userID := utils.GetUserIDFromContext(c)

	result, err := h.groupService.GetGroupInfo(groupID, userID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, result)
}

type CreateGroupRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=50"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
}

func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	userID := utils.GetUserIDFromContext(c)

	result, err := h.groupService.CreateGroup(userID, req.Name, req.Description, req.Avatar)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, result)
}

type JoinGroupRequest struct {
	GroupNumber string `json:"group_number" binding:"required"`
	Message     string `json:"message"`
}

func (h *GroupHandler) JoinGroup(c *gin.Context) {
	var req JoinGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	userID := utils.GetUserIDFromContext(c)

	err := h.groupService.JoinGroup(userID, req.GroupNumber, req.Message)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *GroupHandler) LeaveGroup(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.ParseUint(groupIDStr, 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "群聊ID无效")
		return
	}

	userID := utils.GetUserIDFromContext(c)

	err = h.groupService.LeaveGroup(userID, groupID)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *GroupHandler) UpdateGroupSettings(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.ParseUint(groupIDStr, 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "群聊ID无效")
		return
	}

	var req struct {
		JoinType   int  `json:"join_type"`
		IsAllMuted bool `json:"is_all_muted"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	userID := utils.GetUserIDFromContext(c)

	err = h.groupService.UpdateGroupSettings(groupID, userID, req.JoinType, req.IsAllMuted)
	if err != nil {
		utils.Error(c, http.StatusForbidden, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *GroupHandler) AddAdmin(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.ParseUint(groupIDStr, 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "群聊ID无效")
		return
	}

	var req struct {
		UserID uint64 `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	userID := utils.GetUserIDFromContext(c)

	err = h.groupService.AddAdmin(groupID, userID, req.UserID)
	if err != nil {
		utils.Error(c, http.StatusForbidden, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *GroupHandler) RemoveAdmin(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.ParseUint(groupIDStr, 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "群聊ID无效")
		return
	}

	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "用户ID无效")
		return
	}

	ownerID := utils.GetUserIDFromContext(c)

	err = h.groupService.RemoveAdmin(groupID, ownerID, userID)
	if err != nil {
		utils.Error(c, http.StatusForbidden, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *GroupHandler) ApplyJoinGroup(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.ParseUint(groupIDStr, 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "群聊ID无效")
		return
	}

	userID := utils.GetUserIDFromContext(c)

	err = h.groupService.ApplyJoinGroup(userID, groupID)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *GroupHandler) GetJoinRequests(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)

	requests, err := h.groupService.GetJoinRequests(userID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, requests)
}

func (h *GroupHandler) HandleJoinRequest(c *gin.Context) {
	requestIDStr := c.Param("id")
	requestID, err := strconv.ParseUint(requestIDStr, 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "申请ID无效")
		return
	}

	var req struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	userID := utils.GetUserIDFromContext(c)

	err = h.groupService.HandleJoinRequest(uint(requestID), userID, req.Status)
	if err != nil {
		utils.Error(c, http.StatusForbidden, err.Error())
		return
	}

	utils.Success(c, nil)
}
