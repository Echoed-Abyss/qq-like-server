package service

import (
	"errors"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/echoed-abyss/qq-like-server/internal/model"
	"github.com/echoed-abyss/qq-like-server/pkg/utils"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Nickname string `json:"nickname"`
}

type LoginRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token    string      `json:"token"`
	UserInfo interface{} `json:"user_info"`
}

func (s *AuthService) Register(req *RegisterRequest) (*model.User, error) {
	var count int64
	model.DB.Model(&model.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		return nil, errors.New("用户名已存在")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	qqNumber := generateQQNumber()
	nickname := req.Nickname
	if nickname == "" {
		nickname = req.Username
	}

	user := &model.User{
		QQNumber:     qqNumber,
		Username:     req.Username,
		Nickname:     nickname,
		Password:     string(hashedPassword),
		Avatar:       "",
		Banner:       "",
		Signature:    "这个人很懒，什么都没留下",
		Level:        1,
		LevelExp:     0,
		Gender:       0,
		Age:          0,
		Constellation: "",
		Location:     "",
		Occupation:   "",
		Status:       model.UserStatusNormal,
		OnlineStatus: model.OnlineStatusOffline,
	}

	result := model.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	createDefaultFriendGroups(user.ID)

	return user, nil
}

func (s *AuthService) Login(req *LoginRequest, deviceType int, deviceName string, deviceModel string, ip string) (*LoginResponse, error) {
	var user model.User
	result := model.DB.Where("username = ? OR qq_number = ?", req.Account, req.Account).First(&user)
	if result.Error != nil {
		return nil, errors.New("用户不存在")
	}

	if user.Status != model.UserStatusNormal {
		return nil, errors.New("账号已被禁用")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("密码错误")
	}

	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, errors.New("生成Token失败")
	}

	user.OnlineStatus = model.OnlineStatusOnline
	model.DB.Save(&user)

	s.handleDeviceLogin(user.ID, deviceType, deviceName, deviceModel, ip)

	return &LoginResponse{
		Token:    token,
		UserInfo: user,
	}, nil
}

func (s *AuthService) handleDeviceLogin(userID uint64, deviceType int, deviceName string, deviceModel string, ip string) {
	var devices []model.UserDevice
	model.DB.Where("user_id = ? AND device_type = ? AND is_online = ?", userID, deviceType, true).Find(&devices)

	for _, d := range devices {
		d.IsOnline = false
		model.DB.Save(&d)
	}

	device := model.UserDevice{
		UserID:      userID,
		DeviceType:  deviceType,
		DeviceName:  deviceName,
		DeviceModel: deviceModel,
		IPAddress:   ip,
		LastOnline:  time.Now(),
		IsOnline:    true,
	}
	model.DB.Create(&device)
}

func (s *AuthService) Logout(userID uint64, deviceID string) error {
	var device model.UserDevice
	result := model.DB.Where("user_id = ? AND connection_id = ?", userID, deviceID).First(&device)
	if result.Error == nil {
		device.IsOnline = false
		device.LastOnline = time.Now()
		model.DB.Save(&device)
	}

	var onlineCount int64
	model.DB.Model(&model.UserDevice{}).Where("user_id = ? AND is_online = ?", userID, true).Count(&onlineCount)
	if onlineCount == 0 {
		model.DB.Model(&model.User{}).Where("id = ?", userID).Update("online_status", model.OnlineStatusOffline)
	}

	return nil
}

func generateQQNumber() string {
	rand.Seed(time.Now().UnixNano())
	num := 100000000 + rand.Intn(900000000)
	return int64ToString(int64(num))
}

func createDefaultFriendGroups(userID uint64) {
	groups := []string{"我的好友", "朋友", "家人", "同学", "同事"}
	for i, name := range groups {
		group := model.FriendGroup{
			UserID: userID,
			Name:   name,
			Sort:   i,
		}
		model.DB.Create(&group)
	}
}

func int64ToString(n int64) string {
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
