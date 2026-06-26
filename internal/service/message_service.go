package service

import (
	"errors"
	"time"

	"github.com/echoed-abyss/qq-like-server/internal/model"
)

type MessageService struct{}

func NewMessageService() *MessageService {
	return &MessageService{}
}

type SendMessageRequest struct {
	SessionType int    `json:"session_type"`
	SessionID   uint64 `json:"session_id"`
	MsgType     int    `json:"msg_type"`
	Content     string `json:"content"`
	MediaURL    string `json:"media_url"`
	FileName    string `json:"file_name"`
}

type MessageInfo struct {
	ID            uint64 `json:"id"`
	SessionType   int    `json:"session_type"`
	SessionID     uint64 `json:"session_id"`
	SenderID      uint64 `json:"sender_id"`
	SenderName    string `json:"sender_name"`
	SenderAvatar  string `json:"sender_avatar"`
	MsgType       int    `json:"msg_type"`
	Content       string `json:"content"`
	MediaURL      string `json:"media_url"`
	MediaSize     int64  `json:"media_size"`
	MediaDuration float64 `json:"media_duration"`
	ThumbURL      string `json:"thumb_url"`
	Width         int    `json:"width"`
	Height        int    `json:"height"`
	FileName      string `json:"file_name"`
	IsRecalled    bool   `json:"is_recalled"`
	SendTime      int64  `json:"send_time"`
}

func (s *MessageService) SendMessage(userID uint64, req *SendMessageRequest) (*MessageInfo, error) {
	if req.SessionType != model.SessionTypePrivate && req.SessionType != model.SessionTypeGroup {
		return nil, errors.New("无效的会话类型")
	}

	if req.MsgType < model.MsgTypeText || req.MsgType > model.MsgTypeAt {
		return nil, errors.New("无效的消息类型")
	}

	if req.MsgType == model.MsgTypeText && req.Content == "" {
		return nil, errors.New("消息内容不能为空")
	}

	msg := &model.Message{
		SessionType: req.SessionType,
		SessionID:   req.SessionID,
		SenderID:    userID,
		MsgType:     req.MsgType,
		Content:     req.Content,
		MediaURL:    req.MediaURL,
		FileName:    req.FileName,
		SendTime:    time.Now(),
	}

	result := model.DB.Create(msg)
	if result.Error != nil {
		return nil, result.Error
	}

	return s.messageToInfo(msg), nil
}

func (s *MessageService) RecallMessage(userID uint64, msgID uint64) error {
	var msg model.Message
	result := model.DB.Where("id = ?", msgID).First(&msg)
	if result.Error != nil {
		return errors.New("消息不存在")
	}

	if msg.SenderID != userID {
		return errors.New("只能撤回自己发送的消息")
	}

	if msg.IsRecalled {
		return errors.New("消息已撤回")
	}

	if time.Since(msg.SendTime).Minutes() > 2 {
		return errors.New("超过撤回时间限制")
	}

	msg.IsRecalled = true
	msg.RecalledBy = userID
	msg.RecalledAt = time.Now()
	model.DB.Save(&msg)

	return nil
}

func (s *MessageService) GetMessageList(userID uint64, sessionType int, sessionID uint64, page int, pageSize int) ([]MessageInfo, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	model.DB.Model(&model.Message{}).Where("session_type = ? AND session_id = ?", sessionType, sessionID).Count(&total)

	var messages []model.Message
	offset := (page - 1) * pageSize
	model.DB.Where("session_type = ? AND session_id = ?", sessionType, sessionID).
		Order("send_time DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&messages)

	var result []MessageInfo
	for _, msg := range messages {
		result = append(result, *s.messageToInfo(&msg))
	}

	return result, total, nil
}

func (s *MessageService) AddFavorite(userID uint64, msgID uint64) error {
	var msg model.Message
	result := model.DB.Where("id = ?", msgID).First(&msg)
	if result.Error != nil {
		return errors.New("消息不存在")
	}

	var count int64
	model.DB.Model(&model.Favorite{}).Where("user_id = ? AND msg_id = ?", userID, msgID).Count(&count)
	if count > 0 {
		return errors.New("已收藏该消息")
	}

	favorite := &model.Favorite{
		UserID:    userID,
		MsgID:     msgID,
		MsgType:   msg.MsgType,
		Content:   msg.Content,
		MediaURL:  msg.MediaURL,
		FileName:  msg.FileName,
		SenderID:  msg.SenderID,
		SessionID: msg.SessionID,
	}

	model.DB.Create(favorite)
	return nil
}

func (s *MessageService) RemoveFavorite(userID uint64, favoriteID uint64) error {
	result := model.DB.Where("id = ? AND user_id = ?", favoriteID, userID).Delete(&model.Favorite{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("收藏不存在")
	}
	return nil
}

func (s *MessageService) GetFavorites(userID uint64, page int, pageSize int) ([]model.Favorite, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	model.DB.Model(&model.Favorite{}).Where("user_id = ?", userID).Count(&total)

	var favorites []model.Favorite
	offset := (page - 1) * pageSize
	model.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&favorites)

	return favorites, total, nil
}

func (s *MessageService) messageToInfo(msg *model.Message) *MessageInfo {
	info := &MessageInfo{
		ID:            msg.ID,
		SessionType:   msg.SessionType,
		SessionID:     msg.SessionID,
		SenderID:      msg.SenderID,
		MsgType:       msg.MsgType,
		Content:       msg.Content,
		MediaURL:      msg.MediaURL,
		MediaSize:     msg.MediaSize,
		MediaDuration: msg.MediaDuration,
		ThumbURL:      msg.ThumbURL,
		Width:         msg.Width,
		Height:        msg.Height,
		FileName:      msg.FileName,
		IsRecalled:    msg.IsRecalled,
		SendTime:      msg.SendTime.Unix(),
	}

	var sender model.User
	model.DB.Where("id = ?", msg.SenderID).Select("nickname, avatar").First(&sender)
	info.SenderName = sender.Nickname
	info.SenderAvatar = sender.Avatar

	return info
}
