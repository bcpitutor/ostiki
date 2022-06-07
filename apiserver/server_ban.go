package apiserver

import (
	"github.com/bcpitutor/ostiki/middleware"
	"github.com/bcpitutor/ostiki/routes"
	"github.com/gin-gonic/gin"
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
