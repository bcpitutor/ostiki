package apiserver

import (
	"github.com/gin-gonic/gin"
	"github.com/tiki-systems/tikiserver/middleware"
	"github.com/tiki-systems/tikiserver/routes"
)

func addBanHandlers(ginEngine *gin.Engine, vars middleware.GinHandlerVars) {
	ginEngine.GET("/user/banlist",
		func(ctx *gin.Context) {
			routes.ListBannedUsers(ctx, vars)
		},
	)

	ginEngine.POST("/user/ban",
		func(ctx *gin.Context) {
			routes.BanUser(ctx, vars)
		},
	)

	ginEngine.DELETE("/user/unban",
		func(ctx *gin.Context) {
			routes.UnbanUser(ctx, vars)
		},
	)
}
