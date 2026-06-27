package middleware

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/echoed-abyss/qq-like-server/internal/crypto"
	"github.com/echoed-abyss/qq-like-server/pkg/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Error(c, http.StatusUnauthorized, "未授权访问")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			utils.Error(c, http.StatusUnauthorized, "认证格式错误")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			utils.Error(c, http.StatusUnauthorized, "Token无效或已过期")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Timestamp, X-Nonce, X-Sign")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func SignatureMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		timestampStr := c.GetHeader("X-Timestamp")
		nonce := c.GetHeader("X-Nonce")
		signature := c.GetHeader("X-Sign")

		if timestampStr == "" || nonce == "" || signature == "" {
			utils.Error(c, http.StatusBadRequest, "缺少签名参数")
			c.Abort()
			return
		}

		timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "时间戳格式错误")
			c.Abort()
			return
		}

		var bodyStr string
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			var bodyBytes []byte
			bodyBytes, err = io.ReadAll(c.Request.Body)
			if err != nil {
				utils.Error(c, http.StatusBadRequest, "读取请求体失败")
				c.Abort()
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			bodyStr = string(bodyBytes)
		}

		if !crypto.VerifySignature(bodyStr, timestamp, nonce, signature) {
			utils.Error(c, http.StatusBadRequest, "签名验证失败")
			c.Abort()
			return
		}

		c.Next()
	}
}
