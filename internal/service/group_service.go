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

func (s *GroupService) CreateGroup(userID uint64, name string, description string) (*model.Group, error) {
	groupNumber := generateGroupNumber()

	group := &model.Group{
		GroupNumber: groupNumber,
		Name:        name,
		Description: description,
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
	num := 100000 + rand.Intn(900000)
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
