package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Version godoc
// @Summary Welcome
// @Description welcome message
// @Tags Version
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]any
// @Router / [get]
func Welcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Move on, nothing to see here",
	})
}
