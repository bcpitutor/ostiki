package routes

import (
	"net/http"
	"time"

	"github.com/bcpitutor/ostiki/middleware"
	"github.com/bcpitutor/ostiki/version"
	"github.com/gin-gonic/gin"
)

func VersionInfomation(c *gin.Context, vars middleware.GinHandlerVars) {
	logger := vars.Logger
	logger.Debug("in Version Infomation")

	version := map[string]any{
		"version":      version.VersionDetails.Version,
		"time":         version.VersionDetails.Time,
		"git":          version.VersionDetails.Details,
		"buildMachine": version.VersionDetails.Machine,
	}

	c.JSON(http.StatusOK,
		gin.H{
			"status":  "success",
			"message": "done",
			"details": version,
			"current": time.Now().String(),
		},
	)
}
