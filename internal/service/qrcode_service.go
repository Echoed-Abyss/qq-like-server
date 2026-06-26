package service

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
)

type QRCodeService struct{}

func NewQRCodeService() *QRCodeService {
	return &QRCodeService{}
}

type QRCodeContent struct {
	Type      string `json:"type"`
	ID        uint64 `json:"id"`
	Number    string `json:"number"`
	Timestamp int64  `json:"ts"`
}

const (
	QRTypeUser  = "user"
	QRTypeGroup = "group"
)

func (s *QRCodeService) GenerateUserQRCode(userID uint64, qqNumber string) (string, error) {
	if userID == 0 {
		return "", errors.New("用户ID无效")
	}

	content := QRCodeContent{
		Type:   QRTypeUser,
		ID:     userID,
		Number: qqNumber,
	}

	data, err := json.Marshal(content)
	if err != nil {
		return "", err
	}

	encoded := base64.URLEncoding.EncodeToString(data)
	qrUrl := fmt.Sprintf("qqlike://user/%s", encoded)

	return qrUrl, nil
}

func (s *QRCodeService) GenerateGroupQRCode(groupID uint64, groupNumber string) (string, error) {
	if groupID == 0 {
		return "", errors.New("群聊ID无效")
	}

	content := QRCodeContent{
		Type:   QRTypeGroup,
		ID:     groupID,
		Number: groupNumber,
	}

	data, err := json.Marshal(content)
	if err != nil {
		return "", err
	}

	encoded := base64.URLEncoding.EncodeToString(data)
	qrUrl := fmt.Sprintf("qqlike://group/%s", encoded)

	return qrUrl, nil
}

func (s *QRCodeService) ParseQRCode(qrContent string) (*QRCodeContent, error) {
	var typeStr string
	var encoded string

	_, err := fmt.Sscanf(qrContent, "qqlike://%s/%s", &typeStr, &encoded)
	if err != nil {
		return nil, errors.New("无效的二维码格式")
	}

	data, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, errors.New("二维码解码失败")
	}

	var content QRCodeContent
	err = json.Unmarshal(data, &content)
	if err != nil {
		return nil, errors.New("二维码数据解析失败")
	}

	return &content, nil
}
