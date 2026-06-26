package model

import (
	"time"
)

type User struct {
	ID           uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	QQNumber     string    `json:"qq_number" gorm:"uniqueIndex;size:20"`
	Username     string    `json:"username" gorm:"size:50"`
	Nickname     string    `json:"nickname" gorm:"size:50"`
	Password     string    `json:"-" gorm:"size:255"`
	Avatar       string    `json:"avatar" gorm:"size:500"`
	Banner       string    `json:"banner" gorm:"size:500"`
	Signature    string    `json:"signature" gorm:"size:200"`
	Level        int       `json:"level" gorm:"default:1"`
	LevelExp     int       `json:"level_exp" gorm:"default:0"`
	Gender       int       `json:"gender" gorm:"default:0"`
	Age          int       `json:"age" gorm:"default:0"`
	Constellation string   `json:"constellation" gorm:"size:20"`
	Location     string    `json:"location" gorm:"size:100"`
	Occupation   string    `json:"occupation" gorm:"size:50"`
	Status       int       `json:"status" gorm:"default:1"`
	OnlineStatus int       `json:"online_status" gorm:"default:0"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserDevice struct {
	ID           uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID       uint64    `json:"user_id" gorm:"index"`
	DeviceType   int       `json:"device_type"`
	DeviceName   string    `json:"device_name" gorm:"size:100"`
	DeviceModel  string    `json:"device_model" gorm:"size:100"`
	IPAddress    string    `json:"ip_address" gorm:"size:50"`
	LastOnline   time.Time `json:"last_online"`
	IsOnline     bool      `json:"is_online" gorm:"default:false"`
	ConnectionID string    `json:"connection_id" gorm:"size:100"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Friend struct {
	ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uint64    `json:"user_id" gorm:"index:idx_user_friend,unique"`
	FriendID  uint64    `json:"friend_id" gorm:"index:idx_user_friend,unique"`
	Remark    string    `json:"remark" gorm:"size:50"`
	GroupID   uint64    `json:"group_id"`
	Status    int       `json:"status" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
}

type FriendGroup struct {
	ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uint64    `json:"user_id" gorm:"index"`
	Name      string    `json:"name" gorm:"size:50"`
	Sort      int       `json:"sort" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
}

type Group struct {
	ID          uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	GroupNumber string    `json:"group_number" gorm:"uniqueIndex;size:20"`
	Name        string    `json:"name" gorm:"size:100"`
	Avatar      string    `json:"avatar" gorm:"size:500"`
	Description string    `json:"description" gorm:"size:500"`
	OwnerID     uint64    `json:"owner_id"`
	MemberCount int       `json:"member_count" gorm:"default:0"`
	MaxMembers  int       `json:"max_members" gorm:"default:2000"`
	Status      int       `json:"status" gorm:"default:1"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GroupMember struct {
	ID           uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	GroupID      uint64    `json:"group_id" gorm:"index:idx_group_user,unique"`
	UserID       uint64    `json:"user_id" gorm:"index:idx_group_user,unique"`
	Nickname     string    `json:"nickname" gorm:"size:50"`
	Role         int       `json:"role" gorm:"default:0"`
	JoinTime     time.Time `json:"join_time"`
	LastReadMsgID uint64   `json:"last_read_msg_id" gorm:"default:0"`
	MuteUntil    time.Time `json:"mute_until"`
}

type Message struct {
	ID          uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	SessionType int       `json:"session_type"`
	SessionID   uint64    `json:"session_id" gorm:"index"`
	SenderID    uint64    `json:"sender_id"`
	MsgType     int       `json:"msg_type"`
	Content     string    `json:"content" gorm:"type:text"`
	MediaURL    string    `json:"media_url" gorm:"size:500"`
	MediaSize   int64     `json:"media_size"`
	MediaDuration float64  `json:"media_duration"`
	ThumbURL    string    `json:"thumb_url" gorm:"size:500"`
	Width       int       `json:"width"`
	Height      int       `json:"height"`
	FileName    string    `json:"file_name" gorm:"size:200"`
	IsRecalled  bool      `json:"is_recalled" gorm:"default:false"`
	RecalledBy  uint64    `json:"recalled_by"`
	RecalledAt  time.Time `json:"recalled_at"`
	SendTime    time.Time `json:"send_time"`
	CreatedAt   time.Time `json:"created_at"`
}

type Favorite struct {
	ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uint64    `json:"user_id" gorm:"index"`
	MsgID     uint64    `json:"msg_id"`
	MsgType   int       `json:"msg_type"`
	Content   string    `json:"content" gorm:"type:text"`
	MediaURL  string    `json:"media_url" gorm:"size:500"`
	FileName  string    `json:"file_name" gorm:"size:200"`
	SenderID  uint64    `json:"sender_id"`
	SessionID uint64    `json:"session_id"`
	CreatedAt time.Time `json:"created_at"`
}

type FriendRequest struct {
	ID         uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	FromUserID uint64    `json:"from_user_id" gorm:"index"`
	ToUserID   uint64    `json:"to_user_id" gorm:"index"`
	Message    string    `json:"message" gorm:"size:200"`
	Status     int       `json:"status" gorm:"default:0"`
	HandledAt  time.Time `json:"handled_at"`
	CreatedAt  time.Time `json:"created_at"`
}

type GroupRequest struct {
	ID         uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	GroupID    uint64    `json:"group_id" gorm:"index"`
	UserID     uint64    `json:"user_id" gorm:"index"`
	Message    string    `json:"message" gorm:"size:200"`
	Status     int       `json:"status" gorm:"default:0"`
	HandledBy  uint64    `json:"handled_by"`
	HandledAt  time.Time `json:"handled_at"`
	CreatedAt  time.Time `json:"created_at"`
}

const (
	DeviceTypeAndroid = 1
	DeviceTypeIOS     = 2
	DeviceTypeWindows = 3
	DeviceTypeMac     = 4
	DeviceTypeLinux   = 5
	DeviceTypeWeb     = 6
)

const (
	OnlineStatusOnline   = 0
	OnlineStatusAway     = 1
	OnlineStatusBusy     = 2
	OnlineStatusDoNotDisturb = 3
	OnlineStatusInvisible = 4
	OnlineStatusOffline  = 5
)

const (
	SessionTypePrivate = 1
	SessionTypeGroup   = 2
)

const (
	MsgTypeText     = 1
	MsgTypeImage    = 2
	MsgTypeVideo    = 3
	MsgTypeVoice    = 4
	MsgTypeFile     = 5
	MsgTypeEmoji    = 6
	MsgTypeSystem   = 7
	MsgTypeRecall   = 8
	MsgTypeAt       = 9
)

const (
	UserStatusNormal  = 1
	UserStatusBanned  = 2
	UserStatusDeleted = 3
)

const (
	FriendStatusPending  = 0
	FriendStatusAccepted = 1
	FriendStatusRejected = 2
	FriendStatusDeleted  = 3
)

const (
	GroupRoleMember = 0
	GroupRoleAdmin  = 1
	GroupRoleOwner  = 2
)

const (
	GroupStatusNormal  = 1
	GroupStatusBanned  = 2
	GroupStatusDeleted = 3
)
