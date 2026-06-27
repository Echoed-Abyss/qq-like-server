package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/echoed-abyss/qq-like-server/internal/config"
	"github.com/echoed-abyss/qq-like-server/internal/handler"
	"github.com/echoed-abyss/qq-like-server/internal/middleware"
	"github.com/echoed-abyss/qq-like-server/internal/model"
	"github.com/echoed-abyss/qq-like-server/internal/tcp"
)

func main() {
	config.LoadConfig()

	model.InitDB()

	go func() {
		tcpServer := tcp.GetServer()
		if err := tcpServer.Start(); err != nil {
			log.Fatalf("Failed to start TCP server: %v", err)
		}
	}()

	gin.SetMode(config.AppConfig.Server.Mode)
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.SignatureMiddleware())

	authHandler := handler.NewAuthHandler()
	userHandler := handler.NewUserHandler()
	messageHandler := handler.NewMessageHandler()
	groupHandler := handler.NewGroupHandler()
	qrHandler := handler.NewQRCodeHandler()

	r.GET("/api/health", authHandler.Health)

	api := r.Group("/api/v1")
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
