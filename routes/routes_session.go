package routes

import (
	"fmt"
	"net/http"

	"github.com/bcpitutor/ostiki/middleware"
	"github.com/bcpitutor/ostiki/models"
	"github.com/gin-gonic/gin"
)

func ListSessions(c *gin.Context, vars middleware.GinHandlerVars) {
	sessionRepository := vars.SessionRepository
	groupRepository := vars.GroupRepository

	userEmail := c.Request.Header.Get("email")
	if !groupRepository.IsUserInTikiadmins(userEmail) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("User %s is not authorized to access banned users", userEmail),
			"data":    "",
			"count":   0,
		})
		return
	}

	sessType, ok := c.GetQuery("sessionType")
	if !ok {
		sessType = "all"
	}

	sessions, err := sessionRepository.GetSessions(sessType)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err,
			"data":    "",
			"count":   "",
		})
	}

	var sessExposes []models.SessionExpose

	for _, v := range sessions {
		var sessExpose models.SessionExpose

		sessExpose.SessionId = v.SessID
		sessExpose.SessionOwner = v.SessionOwner
		sessExpose.RefreshCount = v.Rtimes
		sessExpose.ExpiresAt = v.Expire
		sessExpose.SessionDetails = v.Details
		sessExpose.Revoked = v.IsRevoked
		sessExposes = append(sessExposes, sessExpose)
	}

	newToken, _ := c.Get("newToken")

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"message":  "",
		"data":     sessExposes,
		"count":    len(sessExposes),
		"newToken": newToken,
	})
}
