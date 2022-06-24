package routes

import (
	"fmt"
	"net/http"

	"github.com/bcpitutor/ostiki/middleware"
	"github.com/bcpitutor/ostiki/models"
	"github.com/bcpitutor/ostiki/utils"
	"github.com/gin-gonic/gin"
)

func ListSessions(c *gin.Context, vars middleware.GinHandlerVars) {
	sessionRepository := vars.SessionRepository
	groupRepository := vars.GroupRepository
	imoRepository := vars.ImoRepository
	sugar := vars.Logger.Sugar()

	var sessions []models.Session

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

	sessionsFromImo := imoRepository.GetSessions()
	if len(sessionsFromImo) != 0 {
		sugar.Infof("Got sessions in imo: %+v", sessionsFromImo)
		sessions = sessionsFromImo
	} else {
		sugar.Infof("Session are not in imo, reading from DB")
		sessType, ok := c.GetQuery("sessionType")
		if !ok {
			sessType = "all"
		}

		sessionsFromDB, err := sessionRepository.GetSessions(sessType)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": err,
				"data":    "",
				"count":   "",
			})
		}
		sessions = sessionsFromDB
		imoRepository.SetSessions(sessions)
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
		sessExpose.Epoch = v.Epoch
		sessExpose.SessionExpEpoch = v.SessionExpEpoch

		sessExposes = append(sessExposes, sessExpose)
	}

	newToken, _ := c.Get("newToken")
	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"message":  "",
		"data":     sessExposes,
		"myIP":     utils.GetOutboundIP(),
		"count":    len(sessExposes),
		"newToken": newToken,
	})
}
