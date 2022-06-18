package apiserver

import (
	"github.com/bcpitutor/ostiki/middleware"
	"github.com/bcpitutor/ostiki/routes"
	"github.com/gin-gonic/gin"
)

func addMiscHandlers(ginEngine *gin.Engine, vars middleware.GinHandlerVars) {
	ginEngine.GET("/",
		routes.Welcome,
	)

	ginEngine.GET("/version",
		func(ctx *gin.Context) {
			routes.VersionInfomation(ctx, vars)
		},
	)

	ginEngine.GET("/auth/:id/:port",
		func(ctx *gin.Context) {
			routes.InitHandler(ctx, vars)
		},
	)

	ginEngine.GET("/auth/googlecb",
		func(ctx *gin.Context) {
			routes.GoogleCBHandler(ctx, vars)
		},
	)

	ginEngine.POST("/checkToken",
		func(ctx *gin.Context) {
			routes.CheckToken(ctx, vars)
		},
	)

	ginEngine.GET("/session/list",
		func(ctx *gin.Context) {
			routes.ListSessions(ctx, vars)
		},
	)

	ginEngine.Use(
		func(ctx *gin.Context) {
			middleware.Auth(ctx, vars)
		},
	)
}
