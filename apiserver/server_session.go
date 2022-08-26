package apiserver

import (
	"github.com/bcpitutor/ostiki/middleware"
	"github.com/bcpitutor/ostiki/routes"
	"github.com/gin-gonic/gin"
)

func addSessionHandlers(ginEngine *gin.Engine, vars middleware.GinHandlerVars, isCacheReady *bool) {
	ginEngine.GET("/session/list",
		func(ctx *gin.Context) {
			routes.ListSessions(ctx, vars, isCacheReady)
		},
	)
}
