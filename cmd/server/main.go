package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/echoed-abyss/qq-like-server/internal/config"
	"github.com/echoed-abyss/qq-like-server/internal/crypto"
	"github.com/echoed-abyss/qq-like-server/internal/handler"
	"github.com/echoed-abyss/qq-like-server/internal/middleware"
	"github.com/echoed-abyss/qq-like-server/internal/model"
	"github.com/echoed-abyss/qq-like-server/internal/tcp"
)

func main() {
	config.LoadConfig()

	// 初始化加密模块
	crypto.InitCrypto(config.AppConfig.AppSecret)

	model.InitDB()

	go func() {
		tcpServer := tcp.GetServer()
		if err := tcpServer.Start(); err != nil {
			log.Fatalf("Failed to start TCP server: %v", err)
		}
	}()

	gin.SetMode(config.AppConfig.Server.Mode)
	r := gin.Default()

	// 初始化 handler
	authHandler := handler.NewAuthHandler()
	userHandler := handler.NewUserHandler()
	messageHandler := handler.NewMessageHandler()
	groupHandler := handler.NewGroupHandler()
	qrHandler := handler.NewQRCodeHandler()

	// CORS 中间件
	r.Use(middleware.CORSMiddleware())

	// 健康检查不需要签名验证
	r.GET("/api/health", authHandler.Health)

	// 需要签名验证的接口组
	api := r.Group("/api")
	api.Use(middleware.SignatureMiddleware())
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", middleware.AuthMiddleware(), authHandler.Logout)
		}

		user := api.Group("/user")
		user.Use(middleware.AuthMiddleware())
		{
			user.GET("/info", userHandler.GetUserInfo)
			user.GET("/profile/:id", userHandler.GetUserProfile)
			user.PUT("/status", userHandler.UpdateStatus)
			user.GET("/devices", userHandler.GetDevices)
			user.DELETE("/device/:id", userHandler.KickDevice)
		}

		msg := api.Group("/message")
		msg.Use(middleware.AuthMiddleware())
		{
			msg.POST("/send", messageHandler.SendMessage)
			msg.DELETE("/:id/recall", messageHandler.RecallMessage)
			msg.GET("/list", messageHandler.GetMessageList)
			msg.POST("/favorite", messageHandler.AddFavorite)
			msg.DELETE("/favorite/:id", messageHandler.RemoveFavorite)
			msg.GET("/favorites", messageHandler.GetFavorites)
		}

		group := api.Group("/group")
		group.Use(middleware.AuthMiddleware())
		{
			group.GET("/:id", groupHandler.GetGroupInfo)
			group.POST("/create", groupHandler.CreateGroup)
			group.POST("/join", groupHandler.JoinGroup)
			group.DELETE("/:id/leave", groupHandler.LeaveGroup)
		}

		qr := api.Group("/qrcode")
		qr.Use(middleware.AuthMiddleware())
		{
			qr.GET("/user", qrHandler.GetUserQRCode)
			qr.GET("/group/:id", qrHandler.GetGroupQRCode)
			qr.POST("/parse", qrHandler.ParseQRCode)
		}
	}

	addr := fmt.Sprintf("%s:%d", config.AppConfig.Server.Host, config.AppConfig.Server.Port)
	log.Printf("HTTP server listening on %s", addr)
	r.Run(addr)
}
