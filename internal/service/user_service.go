package service

import (
	"errors"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/echoed-abyss/qq-like-server/internal/model"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

type UserProfile struct {
	ID           uint64              `json:"id"`
	QQNumber     string              `json:"qq_number"`
	Nickname     string              `json:"nickname"`
	Avatar       string              `json:"avatar"`
	Banner       string              `json:"banner"`
	Signature    string              `json:"signature"`
	Level        int                 `json:"level"`
	LevelExp     int                 `json:"level_exp"`
	Gender       int                 `json:"gender"`
	Age          int                 `json:"age"`
	Constellation string             `json:"constellation"`
	Location     string              `json:"location"`
	Occupation   string              `json:"occupation"`
	OnlineStatus int                 `json:"online_status"`
}

type UserInfoResponse struct {
	User         UserProfile   `json:"user"`
	Devices      []model.UserDevice `json:"devices"`
	FriendGroups []FriendGroupInfo `json:"friend_groups"`
	GroupList    []GroupInfo   `json:"group_list"`
	FriendList   []FriendInfo  `json:"friend_list"`
}

type FriendGroupInfo struct {
	ID       uint64       `json:"id"`
	Name     string       `json:"name"`
	Sort     int          `json:"sort"`
	Friends  []FriendInfo `json:"friends"`
	OnlineCount int       `json:"online_count"`
	TotalCount int        `json:"total_count"`
}

type FriendInfo struct {
	ID           uint64 `json:"id"`
	QQNumber     string `json:"qq_number"`
	Nickname     string `json:"nickname"`
	Remark       string `json:"remark"`
	Avatar       string `json:"avatar"`
	Signature    string `json:"signature"`
	OnlineStatus int    `json:"online_status"`
}

type GroupInfo struct {
	ID          uint64 `json:"id"`
	GroupNumber string `json:"group_number"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
	MemberCount int    `json:"member_count"`
	OnlineCount int    `json:"online_count"`
	LastMsg     string `json:"last_msg"`
	LastMsgTime string `json:"last_msg_time"`
	UnreadCount int    `json:"unread_count"`
}

func (s *UserService) GetUserInfo(userID uint64) (*UserInfoResponse, error) {
	var user model.User
	result := model.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return nil, errors.New("用户不存在")
	}

	var devices []model.UserDevice
	model.DB.Where("user_id = ?", userID).Order("is_online DESC, last_online DESC").Find(&devices)

	friendGroups := s.getFriendGroups(userID)
	groupList := s.getUserGroups(userID)
	friendList := s.getFriends(userID)

	return &UserInfoResponse{
		User: UserProfile{
			ID:           user.ID,
			QQNumber:     user.QQNumber,
			Nickname:     user.Nickname,
			Avatar:       user.Avatar,
			Banner:       user.Banner,
			Signature:    user.Signature,
			Level:        user.Level,
			LevelExp:     user.LevelExp,
			Gender:       user.Gender,
			Age:          user.Age,
			Constellation: user.Constellation,
			Location:     user.Location,
			Occupation:   user.Occupation,
			OnlineStatus: user.OnlineStatus,
		},
		Devices:      devices,
		FriendGroups: friendGroups,
		GroupList:    groupList,
		FriendList:   friendList,
	}, nil
}

func (s *UserService) GetUserProfile(targetID uint64, currentUserID uint64) (*UserProfile, error) {
	var user model.User
	result := model.DB.Where("id = ?", targetID).First(&user)
	if result.Error != nil {
		return nil, errors.New("用户不存在")
	}

	return &UserProfile{
		ID:           user.ID,
		QQNumber:     user.QQNumber,
		Nickname:     user.Nickname,
		Avatar:       user.Avatar,
		Banner:       user.Banner,
		Signature:    user.Signature,
		Level:        user.Level,
		LevelExp:     user.LevelExp,
		Gender:       user.Gender,
		Age:          user.Age,
		Constellation: user.Constellation,
		Location:     user.Location,
		Occupation:   user.Occupation,
		OnlineStatus: user.OnlineStatus,
	}, nil
}

func (s *UserService) UpdateOnlineStatus(userID uint64, status int) error {
	result := model.DB.Model(&model.User{}).Where("id = ?", userID).Update("online_status", status)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *UserService) GetUserDevices(userID uint64) ([]model.UserDevice, error) {
	var devices []model.UserDevice
	result := model.DB.Where("user_id = ?", userID).Order("is_online DESC, last_online DESC").Find(&devices)
	if result.Error != nil {
		return nil, result.Error
	}
	return devices, nil
}

func (s *UserService) KickDevice(userID uint64, deviceID uint64) error {
	var device model.UserDevice
	result := model.DB.Where("id = ? AND user_id = ?", deviceID, userID).First(&device)
	if result.Error != nil {
		return errors.New("设备不存在")
	}

	device.IsOnline = false
	model.DB.Save(&device)

	return nil
}

func (s *UserService) getFriendGroups(userID uint64) []FriendGroupInfo {
	var groups []model.FriendGroup
	model.DB.Where("user_id = ?", userID).Order("sort ASC").Find(&groups)

	var result []FriendGroupInfo
	for _, g := range groups {
		friends := s.getFriendsByGroup(userID, g.ID)
		onlineCount := 0
		for _, f := range friends {
			if f.OnlineStatus != model.OnlineStatusOffline {
				onlineCount++
			}
		}
		result = append(result, FriendGroupInfo{
			ID:          g.ID,
			Name:        g.Name,
			Sort:        g.Sort,
			Friends:     friends,
			OnlineCount: onlineCount,
			TotalCount:  len(friends),
		})
	}
	return result
}

func (s *UserService) getFriendsByGroup(userID uint64, groupID uint64) []FriendInfo {
	var friends []model.Friend
	model.DB.Where("user_id = ? AND group_id = ? AND status = ?", userID, groupID, model.FriendStatusAccepted).Find(&friends)

	var result []FriendInfo
	for _, f := range friends {
		friendUser := s.getUserByID(f.FriendID)
		if friendUser != nil {
			result = append(result, FriendInfo{
				ID:           friendUser.ID,
				QQNumber:     friendUser.QQNumber,
				Nickname:     friendUser.Nickname,
				Remark:       f.Remark,
				Avatar:       friendUser.Avatar,
				Signature:    friendUser.Signature,
				OnlineStatus: friendUser.OnlineStatus,
			})
		}
	}
	return result
}

func (s *UserService) getFriends(userID uint64) []FriendInfo {
	var friends []model.Friend
	model.DB.Where("user_id = ? AND status = ?", userID, model.FriendStatusAccepted).Find(&friends)

	var result []FriendInfo
	for _, f := range friends {
		friendUser := s.getUserByID(f.FriendID)
		if friendUser != nil {
			result = append(result, FriendInfo{
				ID:           friendUser.ID,
				QQNumber:     friendUser.QQNumber,
				Nickname:     friendUser.Nickname,
				Remark:       f.Remark,
				Avatar:       friendUser.Avatar,
				Signature:    friendUser.Signature,
				OnlineStatus: friendUser.OnlineStatus,
			})
		}
	}
	return result
}

func (s *UserService) getUserGroups(userID uint64) []GroupInfo {
	var groupMembers []model.GroupMember
	model.DB.Where("user_id = ?", userID).Find(&groupMembers)

	var result []GroupInfo
	addedGroupIDs := make(map[uint64]bool)

	for _, gm := range groupMembers {
		var group model.Group
		model.DB.Where("id = ?", gm.GroupID).First(&group)
		if group.ID > 0 {
			result = append(result, GroupInfo{
				ID:          group.ID,
				GroupNumber: group.GroupNumber,
				Name:        group.Name,
				Avatar:      group.Avatar,
				Description: group.Description,
				MemberCount: group.MemberCount,
				OnlineCount: 0,
				LastMsg:     "",
				LastMsgTime: "",
				UnreadCount: 0,
			})
			addedGroupIDs[group.ID] = true
		}
	}

	var officialGroup model.Group
	model.DB.Where("group_number = ?", "2182344375").First(&officialGroup)
	if officialGroup.ID > 0 && !addedGroupIDs[officialGroup.ID] {
		result = append([]GroupInfo{{
			ID:          officialGroup.ID,
			GroupNumber: officialGroup.GroupNumber,
			Name:        officialGroup.Name,
			Avatar:      officialGroup.Avatar,
			Description: officialGroup.Description,
			MemberCount: officialGroup.MemberCount,
			OnlineCount: 0,
			LastMsg:     "",
			LastMsgTime: "",
			UnreadCount: 0,
		}}, result...)
	}

	return result
}

func (s *UserService) getUserByID(userID uint64) *model.User {
	var user model.User
	result := model.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return nil
	}
	return &user
}

func CalculateLevel(exp int) int {
	if exp < 100 {
		return 1
	}
	levels := []struct{ maxLevel, expPerLevel int }{
		{10, 100}, {20, 150}, {30, 200}, {40, 300}, {50, 500},
	}
	remainingExp := exp
	currentLevel := 1
	for _, tier := range levels {
		for currentLevel < tier.maxLevel && remainingExp >= tier.expPerLevel {
			remainingExp -= tier.expPerLevel
			currentLevel++
		}
		if currentLevel >= tier.maxLevel {
			continue
		}
		break
	}
	if currentLevel > 50 {
		currentLevel = 50
	}
	return currentLevel
}

func (s *UserService) CheckIn(userID uint64) (int, int, error) {
	var user model.User
	result := model.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return 0, 0, errors.New("用户不存在")
	}

	now := time.Now()
	if user.LastCheckIn.Year() == now.Year() && user.LastCheckIn.YearDay() == now.YearDay() {
		return 0, 0, errors.New("今天已经打卡过了")
	}

	expGain := 20 + rand.Intn(41)
	user.Exp += expGain
	user.LastCheckIn = now
	user.Level = CalculateLevel(user.Exp)

	model.DB.Save(&user)
	return expGain, user.Level, nil
}

func (s *UserService) LikeUser(targetID uint64) (int, error) {
	var user model.User
	result := model.DB.Where("id = ?", targetID).First(&user)
	if result.Error != nil {
		return 0, errors.New("用户不存在")
	}

	user.Likes++
	model.DB.Save(&user)
	return user.Likes, nil
}

func (s *UserService) GetLikeRank() ([]model.User, error) {
	var users []model.User
	result := model.DB.Order("likes DESC").Limit(50).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (s *UserService) UpdateProfile(userID uint64, nickname, signature, bio, tags string, gender, age int) error {
	result := model.DB.Model(&model.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"nickname":  nickname,
		"signature": signature,
		"bio":       bio,
		"tags":      tags,
		"gender":    gender,
		"age":       age,
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *UserService) UpdatePassword(userID uint64, oldPassword, newPassword string) error {
	var user model.User
	result := model.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return errors.New("用户不存在")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("原密码错误")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("密码加密失败")
	}

	model.DB.Model(&user).Update("password", string(hashedPassword))
	return nil
}

func (s *UserService) DeleteAccount(userID uint64) error {
	var user model.User
	result := model.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return errors.New("用户不存在")
	}

	model.DB.Delete(&user)
	return nil
}
