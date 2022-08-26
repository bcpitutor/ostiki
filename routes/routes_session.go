package routes

import (
	"fmt"
	"net/http"

	"github.com/bcpitutor/ostiki/middleware"
	"github.com/bcpitutor/ostiki/models"
	"github.com/bcpitutor/ostiki/utils"
	"github.com/gin-gonic/gin"
)

func ListSessions(c *gin.Context, vars middleware.GinHandlerVars, isCacheReady *bool) {
	sessionRepository := vars.SessionRepository
	groupRepository := vars.GroupRepository
	//imoRepository := vars.ImoRepository
	sugar := vars.Logger.Sugar()

	userEmail := c.Request.Header.Get("email")
	//sugar.Infof("User email: %s", userEmail)
	if !groupRepository.IsUserInTikiadmins(userEmail) {
		sugar.Infof("User is not in Tikiadmins group, aborting")
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
		sugar.Infof("No session type specified, reading all sessions")
		sessType = "all"
	}

	//sugar.Info("We are in ListSessions with isCacheReady: ", *isCacheReady)
	var sessions []models.Session

	sessionsFromDB, err := sessionRepository.GetSessions(sessType)
	if err != nil {
		sugar.Errorf("Error getting sessions from sessionRepository: %+v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err,
			"data":    "",
			"count":   "",
		})
	}
	sessions = sessionsFromDB

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
	//sugar.Infof("Returning sessions: %+v", sessExposes)

	respond := gin.H{
		"status":   "success",
		"message":  "",
		"data":     sessExposes,
		"myIP":     utils.GetOutboundIP(),
		"count":    len(sessExposes),
		"newToken": newToken,
	}

	//sugar.Infof("Returning sessions: %+v", respond)

	c.JSON(http.StatusOK, respond)
}
