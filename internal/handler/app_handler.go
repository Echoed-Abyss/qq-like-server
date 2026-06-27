package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/echoed-abyss/qq-like-server/pkg/utils"
)

type AppHandler struct{}

func NewAppHandler() *AppHandler {
	return &AppHandler{}
}

func (h *AppHandler) CheckUpdate(c *gin.Context) {
	version := c.Query("version")
	if version == "" {
		version = "0"
	}

	currentVersion := "1.0.0"
	hasUpdate := false
	updateUrl := ""
	updateLog := "当前已是最新版本"

	if version != currentVersion {
		hasUpdate = false
	}

	utils.Success(c, gin.H{
		"has_update":  hasUpdate,
		"version":     currentVersion,
		"update_url":  updateUrl,
		"update_log":  updateLog,
		"force":       false,
	})
}
