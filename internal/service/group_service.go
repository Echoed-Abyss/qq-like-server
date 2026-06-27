package service

import (
	"errors"
	"math/rand"
	"time"

	"github.com/echoed-abyss/qq-like-server/internal/model"
)

type GroupService struct{}

func NewGroupService() *GroupService {
	return &GroupService{}
}

type GroupDetail struct {
	ID            uint64            `json:"id"`
	GroupNumber   string            `json:"group_number"`
	Name          string            `json:"name"`
	Avatar        string            `json:"avatar"`
	Description   string            `json:"description"`
	OwnerID       uint64            `json:"owner_id"`
	OwnerName     string            `json:"owner_name"`
	MemberCount   int               `json:"member_count"`
	MaxMembers    int               `json:"max_members"`
	OnlineCount   int               `json:"online_count"`
	MemberList    []GroupMemberInfo `json:"member_list"`
	IsMember      bool              `json:"is_member"`
	MyRole        int               `json:"my_role"`
}

type GroupMemberInfo struct {
	UserID       uint64 `json:"user_id"`
	Nickname     string `json:"nickname"`
	GroupNickname string `json:"group_nickname"`
	Avatar       string `json:"avatar"`
	Role         int    `json:"role"`
	JoinTime     int64  `json:"join_time"`
	OnlineStatus int    `json:"online_status"`
}

func (s *GroupService) GetGroupInfo(groupID uint64, userID uint64) (*GroupDetail, error) {
	var group model.Group
	result := model.DB.Where("id = ?", groupID).First(&group)
	if result.Error != nil {
		return nil, errors.New("群聊不存在")
	}

	var myMember model.GroupMember
	model.DB.Where("group_id = ? AND user_id = ?", groupID, userID).First(&myMember)
	isMember := myMember.ID > 0

	var members []model.GroupMember
	model.DB.Where("group_id = ?", groupID).Find(&members)

	var memberList []GroupMemberInfo
	onlineCount := 0
	for _, m := range members {
		var user model.User
		model.DB.Where("id = ?", m.UserID).Select("nickname, avatar, online_status").First(&user)
		memberList = append(memberList, GroupMemberInfo{
			UserID:       m.UserID,
			Nickname:     user.Nickname,
			GroupNickname: m.Nickname,
			Avatar:       user.Avatar,
			Role:         m.Role,
			JoinTime:     m.JoinTime.Unix(),
			OnlineStatus: user.OnlineStatus,
		})
		if user.OnlineStatus != model.OnlineStatusOffline {
			onlineCount++
		}
	}

	var owner model.User
	model.DB.Where("id = ?", group.OwnerID).Select("nickname").First(&owner)

	return &GroupDetail{
		ID:          group.ID,
		GroupNumber: group.GroupNumber,
		Name:        group.Name,
		Avatar:      group.Avatar,
		Description: group.Description,
		OwnerID:     group.OwnerID,
		OwnerName:   owner.Nickname,
		MemberCount: group.MemberCount,
		MaxMembers:  group.MaxMembers,
		OnlineCount: onlineCount,
		MemberList:  memberList,
		IsMember:    isMember,
		MyRole:      myMember.Role,
	}, nil
}

func (s *GroupService) CreateGroup(userID uint64, name string, description string, avatar string) (*model.Group, error) {
	groupNumber := generateGroupNumber()

	group := &model.Group{
		GroupNumber: groupNumber,
		Name:        name,
		Description: description,
		Avatar:      avatar,
		OwnerID:     userID,
		MemberCount: 1,
		MaxMembers:  2000,
		Status:      model.GroupStatusNormal,
	}

	result := model.DB.Create(group)
	if result.Error != nil {
		return nil, result.Error
	}

	member := &model.GroupMember{
		GroupID:  group.ID,
		UserID:   userID,
		Nickname: "",
		Role:     model.GroupRoleOwner,
		JoinTime: time.Now(),
	}
	model.DB.Create(member)

	return group, nil
}

func (s *GroupService) JoinGroup(userID uint64, groupNumber string, message string) error {
	var group model.Group
	result := model.DB.Where("group_number = ?", groupNumber).First(&group)
	if result.Error != nil {
		return errors.New("群聊不存在")
	}

	var count int64
	model.DB.Model(&model.GroupMember{}).Where("group_id = ? AND user_id = ?", group.ID, userID).Count(&count)
	if count > 0 {
		return errors.New("已在群聊中")
	}

	if group.MemberCount >= group.MaxMembers {
		return errors.New("群聊人数已满")
	}

	request := &model.GroupRequest{
		GroupID: group.ID,
		UserID:  userID,
		Message: message,
		Status:  0,
	}
	model.DB.Create(request)

	return nil
}

func (s *GroupService) LeaveGroup(userID uint64, groupID uint64) error {
	var member model.GroupMember
	result := model.DB.Where("group_id = ? AND user_id = ?", groupID, userID).First(&member)
	if result.Error != nil {
		return errors.New("不在群聊中")
	}

	if member.Role == model.GroupRoleOwner {
		return errors.New("群主不能退群，请先转让群主")
	}

	model.DB.Delete(&member)

	model.DB.Model(&model.Group{}).Where("id = ?", groupID).UpdateColumn("member_count", model.DB.Model(&model.GroupMember{}).Where("group_id = ?", groupID).RowsAffected)

	return nil
}

func generateGroupNumber() string {
	rand.Seed(time.Now().UnixNano())
	length := 8 + rand.Intn(5)
	min := 1
	for i := 1; i < length; i++ {
		min *= 10
	}
	max := min * 10
	num := min + rand.Intn(max-min)
	return int64ToStr(int64(num))
}

func int64ToStr(n int64) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	neg := n < 0
	if neg {
		n = -n
	}
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

func (s *GroupService) UpdateGroupSettings(groupID uint64, userID uint64, joinType int, isAllMuted bool) error {
	var group model.Group
	result := model.DB.Where("id = ?", groupID).First(&group)
	if result.Error != nil {
		return errors.New("群聊不存在")
	}

	if group.OwnerID != userID {
		return errors.New("没有权限")
	}

	model.DB.Model(&group).Updates(map[string]interface{}{
		"join_type":    joinType,
		"is_all_muted": isAllMuted,
	})
	return nil
}

func (s *GroupService) AddAdmin(groupID uint64, ownerID uint64, userID uint64) error {
	var group model.Group
	result := model.DB.Where("id = ?", groupID).First(&group)
	if result.Error != nil {
		return errors.New("群聊不存在")
	}

	if group.OwnerID != ownerID {
		return errors.New("没有权限")
	}

	var member model.GroupMember
	result = model.DB.Where("group_id = ? AND user_id = ?", groupID, userID).First(&member)
	if result.Error != nil {
		return errors.New("用户不在群聊中")
	}

	member.Role = model.GroupRoleAdmin
	model.DB.Save(&member)

	var count int64
	model.DB.Model(&model.GroupAdmin{}).Where("group_id = ? AND user_id = ?", groupID, userID).Count(&count)
	if count == 0 {
		admin := model.GroupAdmin{
			GroupID: uint(groupID),
			UserID:  uint(userID),
		}
		model.DB.Create(&admin)
	}

	return nil
}

func (s *GroupService) RemoveAdmin(groupID uint64, ownerID uint64, userID uint64) error {
	var group model.Group
	result := model.DB.Where("id = ?", groupID).First(&group)
	if result.Error != nil {
		return errors.New("群聊不存在")
	}

	if group.OwnerID != ownerID {
		return errors.New("没有权限")
	}

	var member model.GroupMember
	result = model.DB.Where("group_id = ? AND user_id = ?", groupID, userID).First(&member)
	if result.Error != nil {
		return errors.New("用户不在群聊中")
	}

	member.Role = model.GroupRoleMember
	model.DB.Save(&member)

	model.DB.Where("group_id = ? AND user_id = ?", groupID, userID).Delete(&model.GroupAdmin{})
	return nil
}

func (s *GroupService) ApplyJoinGroup(userID uint64, groupID uint64) error {
	var group model.Group
	result := model.DB.Where("id = ?", groupID).First(&group)
	if result.Error != nil {
		return errors.New("群聊不存在")
	}

	var count int64
	model.DB.Model(&model.GroupMember{}).Where("group_id = ? AND user_id = ?", groupID, userID).Count(&count)
	if count > 0 {
		return errors.New("已在群聊中")
	}

	if group.MemberCount >= group.MaxMembers {
		return errors.New("群聊人数已满")
	}

	if group.JoinType == model.GroupJoinTypeForbidden {
		return errors.New("该群不允许加入")
	}

	if group.JoinType == model.GroupJoinTypeAnyone {
		member := &model.GroupMember{
			GroupID:  groupID,
			UserID:   userID,
			Nickname: "",
			Role:     model.GroupRoleMember,
			JoinTime: time.Now(),
		}
		model.DB.Create(member)
		model.DB.Model(&group).UpdateColumn("member_count", model.DB.Model(&model.GroupMember{}).Where("group_id = ?", groupID).RowsAffected)
		return nil
	}

	request := &model.GroupJoinRequest{
		GroupID: uint(groupID),
		UserID:  uint(userID),
		Status:  model.GroupJoinRequestPending,
	}
	model.DB.Create(request)
	return nil
}

func (s *GroupService) GetJoinRequests(userID uint64) ([]model.GroupJoinRequest, error) {
	var groups []model.Group
	model.DB.Where("owner_id = ?", userID).Find(&groups)

	var groupIDs []uint
	for _, g := range groups {
		groupIDs = append(groupIDs, uint(g.ID))
	}

	var adminGroups []model.GroupAdmin
	model.DB.Where("user_id = ?", userID).Find(&adminGroups)
	for _, ag := range adminGroups {
		var g model.Group
		model.DB.Where("id = ?", ag.GroupID).First(&g)
		if g.ID > 0 {
			found := false
			for _, id := range groupIDs {
				if id == uint(g.ID) {
					found = true
					break
				}
			}
			if !found {
				groupIDs = append(groupIDs, uint(g.ID))
			}
		}
	}

	var requests []model.GroupJoinRequest
	if len(groupIDs) > 0 {
		model.DB.Where("group_id IN ? AND status = ?", groupIDs, model.GroupJoinRequestPending).Find(&requests)
	}
	return requests, nil
}

func (s *GroupService) HandleJoinRequest(requestID uint, handlerID uint64, status int) error {
	var req model.GroupJoinRequest
	result := model.DB.Where("id = ?", requestID).First(&req)
	if result.Error != nil {
		return errors.New("申请不存在")
	}

	var group model.Group
	result = model.DB.Where("id = ?", req.GroupID).First(&group)
	if result.Error != nil {
		return errors.New("群聊不存在")
	}

	isOwner := group.OwnerID == handlerID
	var isAdmin bool
	if !isOwner {
		var admin model.GroupAdmin
		model.DB.Where("group_id = ? AND user_id = ?", req.GroupID, handlerID).First(&admin)
		isAdmin = admin.ID > 0
	}

	if !isOwner && !isAdmin {
		return errors.New("没有权限")
	}

	req.Status = status
	model.DB.Save(&req)

	if status == model.GroupJoinRequestAccepted {
		var count int64
		model.DB.Model(&model.GroupMember{}).Where("group_id = ? AND user_id = ?", req.GroupID, req.UserID).Count(&count)
		if count == 0 {
			member := &model.GroupMember{
				GroupID:  uint64(req.GroupID),
				UserID:   uint64(req.UserID),
				Nickname: "",
				Role:     model.GroupRoleMember,
				JoinTime: time.Now(),
			}
			model.DB.Create(member)
			model.DB.Model(&model.Group{}).Where("id = ?", req.GroupID).UpdateColumn("member_count", model.DB.Model(&model.GroupMember{}).Where("group_id = ?", req.GroupID).RowsAffected)
		}
	}

	return nil
}
