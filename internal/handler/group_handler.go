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
}

func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	userID := utils.GetUserIDFromContext(c)

	result, err := h.groupService.CreateGroup(userID, req.Name, req.Description)
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
