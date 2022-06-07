package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/bcpitutor/ostiki/middleware"
	"github.com/bcpitutor/ostiki/models"
	"github.com/gin-gonic/gin"
)

func ListBannedUsers(c *gin.Context, vars middleware.GinHandlerVars) {
	banRepository := vars.BanRepository
	groupRepository := vars.GroupRepository

	adminEmail := c.Request.Header.Get("email")
	if !groupRepository.IsUserInTikiadmins(adminEmail) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("User %s is not authorized to access banned users", adminEmail),
			"data":    "",
			"count":   "",
		})
		return
	}

	resp, err := banRepository.GetBannedUsers()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err,
			"data":    "",
			"count":   "",
		})
		return
	}

	newToken, _ := c.Get("newToken")
	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"message":  "",
		"data":     resp,
		"count":    len(resp),
		"newToken": newToken,
	})
}

func BanUser(c *gin.Context, vars middleware.GinHandlerVars) {
	banRepository := vars.BanRepository
	groupRepository := vars.GroupRepository

	adminEmail := c.Request.Header.Get("email")
	if !groupRepository.IsUserInTikiadmins(adminEmail) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("User %s is not authorized to access banned users", adminEmail),
			"data":    "",
			"count":   "",
		})
		return
	}

	userEmail, ok := c.GetQuery("userEmail")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "please check your query parameters",
		})
		return
	}

	bannedUser := models.BannedUser{
		UserEmail: userEmail,
		CreatedAt: strconv.FormatInt(time.Now().Unix(), 10),
		CreatedBy: adminEmail,
		UpdatedAt: strconv.FormatInt(time.Now().Unix(), 10),
		UpdatedBy: adminEmail,
	}

	err := banRepository.AddBannedUser(bannedUser)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err,
		})
		return
	}

	newToken, _ := c.Get("newToken")

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"message":  fmt.Sprintf("The user, [%s] is banned.", userEmail),
		"newToken": newToken,
	})
}

func UnbanUser(c *gin.Context, vars middleware.GinHandlerVars) {
	banRepository := vars.BanRepository
	groupRepository := vars.GroupRepository

	adminEmail := c.Request.Header.Get("email")
	if !groupRepository.IsUserInTikiadmins(adminEmail) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("User %s is not authorized to access banned users", adminEmail),
			"data":    "",
			"count":   "",
		})
		return
	}

	userEmail, ok := c.GetQuery("userEmail")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "please check your query parameters",
		})
		return
	}

	err := banRepository.UnbanUser(userEmail)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err,
		})
		return
	}

	newToken, _ := c.Get("newToken")

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"message":  fmt.Sprintf("The user, [%s] is unbanned.", userEmail),
		"newToken": newToken,
	})
}
