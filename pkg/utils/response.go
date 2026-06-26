package utils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/echoed-abyss/qq-like-server/internal/config"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Sign    string      `json:"sign,omitempty"`
	Time    int64       `json:"time"`
	Nonce   string      `json:"nonce,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	resp := Response{
		Code:    0,
		Message: "success",
		Data:    data,
		Time:    time.Now().Unix(),
	}
	c.JSON(http.StatusOK, resp)
}

func Error(c *gin.Context, code int, message string) {
	resp := Response{
		Code:    code,
		Message: message,
		Data:    nil,
		Time:    time.Now().Unix(),
	}
	c.JSON(http.StatusOK, resp)
}

func ErrorWithData(c *gin.Context, code int, message string, data interface{}) {
	resp := Response{
		Code:    code,
		Message: message,
		Data:    data,
		Time:    time.Now().Unix(),
	}
	c.JSON(http.StatusOK, resp)
}

type Claims struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint64, username string) (string, error) {
	cfg := config.AppConfig.Token
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.ExpireHour) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "qq-like-server",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

func ParseToken(tokenString string) (*Claims, error) {
	cfg := config.AppConfig.Token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

func GetUserIDFromContext(c *gin.Context) uint64 {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	return userID.(uint64)
}

func GenerateDeviceID() string {
	return fmt.Sprintf("dev_%d", time.Now().UnixNano())
}

func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}
