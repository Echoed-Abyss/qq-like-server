package tcp

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/echoed-abyss/qq-like-server/internal/config"
	"github.com/echoed-abyss/qq-like-server/internal/model"
	"github.com/echoed-abyss/qq-like-server/pkg/utils"
)

type TCPServer struct {
	listener net.Listener
	clients  map[string]*TCPClient
	mu       sync.RWMutex
}

type TCPClient struct {
	Conn       net.Conn
	UserID     uint64
	DeviceID   string
	DeviceType int
	DeviceName string
	LastPing   time.Time
	SendChan   chan *PushMessage
}

type PushMessage struct {
	Type      string      `json:"type"`
	Content   interface{} `json:"content"`
	MsgID     uint64      `json:"msg_id"`
	Timestamp int64       `json:"timestamp"`
}

const (
	MsgTypeNewMessage = "new_message"
	MsgTypeRecall     = "recall_message"
	MsgTypeSystem     = "system"
	MsgTypeKickDevice   = "kick_device"
	MsgTypeFriendOnline  = "friend_online"
	MsgTypeFriendOffline = "friend_offline"
	MsgTypePing      = "ping"
	MsgTypePong      = "pong"
	MsgTypeHeartbeat = "heartbeat"
	MsgTypeLogin     = "login"
)

var server *TCPServer

func NewTCPServer() *TCPServer {
	return &TCPServer{
		clients: make(map[string]*TCPClient),
	}
}

func (s *TCPServer) Start() error {
	cfg := config.AppConfig.TCP
	addr := cfg.Host + ":" + fmt.Sprintf("%d", cfg.Port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.listener = listener
	log.Printf("TCP server listening on %s", addr)

	go s.acceptLoop()
	return nil
}

func (s *TCPServer) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("TCP accept error: %v", err)
			continue
		}

		client := &TCPClient{
			Conn:     conn,
			SendChan: make(chan *PushMessage, 100),
			LastPing: time.Now(),
		}

		go s.handleClient(client)
	}
}

func (s *TCPServer) handleClient(client *TCPClient) {
	defer func() {
		s.removeClient(client)
		client.Conn.Close()
		close(client.SendChan)
	}()

	go s.sendLoop(client)

	buf := make([]byte, 4096)
	for {
		n, err := client.Conn.Read(buf)
		if err != nil {
			log.Printf("TCP read error: %v", err)
			return
		}

		s.handleMessage(client, buf[:n])
	}
}

func (s *TCPServer) handleMessage(client *TCPClient, data []byte) {
	var msg struct {
		Type    string      `json:"type"`
		Token   string      `json:"token"`
		Content interface{} `json:"content"`
		DeviceID string    `json:"device_id"`
		DeviceType int     `json:"device_type"`
		DeviceName string    `json:"device_name"`
	}

	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.Printf("TCP parse error: %v", err)
		return
	}

	switch msg.Type {
	case MsgTypeLogin:
		s.handleLogin(client, msg.Token, msg.DeviceID, msg.DeviceType, msg.DeviceName)
	case MsgTypePing:
		client.LastPing = time.Now()
		s.sendMessage(client, &PushMessage{
			Type:      MsgTypePong,
			Timestamp: time.Now().Unix(),
		})
	case MsgTypeHeartbeat:
		client.LastPing = time.Now()
	}
}

func (s *TCPServer) handleLogin(client *TCPClient, token string, deviceID string, deviceType int, deviceName string) {
	claims, err := utils.ParseToken(token)
	if err != nil {
		s.sendMessage(client, &PushMessage{
			Type:      "login_fail",
			Content:   "无效的token",
			Timestamp: time.Now().Unix(),
		})
		return
	}

	client.UserID = claims.UserID
	client.DeviceID = deviceID
	client.DeviceType = deviceType
	client.DeviceName = deviceName

	s.addClient(client)

	model.DB.Model(&model.UserDevice{}).Where("connection_id = ?", deviceID).Updates(map[string]interface{}{
		"is_online": true,
		"last_online": time.Now(),
	})

	model.DB.Model(&model.User{}).Where("id = ?", claims.UserID).Update("online_status", model.OnlineStatusOnline)

	s.sendMessage(client, &PushMessage{
		Type:      "login_success",
		Content:   "登录成功",
		Timestamp: time.Now().Unix(),
	})

	log.Printf("User %d connected via TCP", claims.UserID)
}

func (s *TCPServer) sendLoop(client *TCPClient) {
	for msg := range client.SendChan {
		data, err := json.Marshal(msg)
		if err != nil {
			continue
		}
		client.Conn.Write(data)
	}
}

func (s *TCPServer) sendMessage(client *TCPClient, msg *PushMessage) {
	select {
	case client.SendChan <- msg:
	default:
		log.Printf("Send channel full for user %d", client.UserID)
	}
}

func (s *TCPServer) addClient(client *TCPClient) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[client.DeviceID] = client
}

func (s *TCPServer) removeClient(client *TCPClient) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, client.DeviceID)

	if client.UserID > 0 {
		model.DB.Model(&model.UserDevice{}).Where("connection_id = ?", client.DeviceID).Update("is_online", false)
	}
}

func (s *TCPServer) PushToUser(userID uint64, msg *PushMessage) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, client := range s.clients {
		if client.UserID == userID {
			s.sendMessage(client, msg)
		}
	}
}

func (s *TCPServer) PushToGroup(groupID uint64, msg *PushMessage) {
	var members []model.GroupMember
	model.DB.Where("group_id = ?", groupID).Find(&members)

	for _, m := range members {
		s.PushToUser(m.UserID, msg)
	}
}

func (s *TCPServer) KickDevice(userID uint64, deviceID string) {
	s.mu.RLock()
	client, exists := s.clients[deviceID]
	s.mu.RUnlock()

	if exists && client.UserID == userID {
		s.sendMessage(client, &PushMessage{
			Type:      MsgTypeKickDevice,
			Content:   "您的账号已在其他设备登录",
			Timestamp: time.Now().Unix(),
		})
		client.Conn.Close()
	}
}

func GetServer() *TCPServer {
	if server == nil {
		server = NewTCPServer()
	}
	return server
}
