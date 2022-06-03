package apiserver

import (
	"github.com/gin-gonic/gin"
	"github.com/tiki-systems/tikiserver/middleware"
	"github.com/tiki-systems/tikiserver/routes"
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

	ginEngine.Use(
		func(ctx *gin.Context) {
			middleware.Auth(ctx, vars)
		},
	)
}
