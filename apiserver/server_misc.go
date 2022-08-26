package apiserver

import (
	"github.com/bcpitutor/ostiki/middleware"
	"github.com/bcpitutor/ostiki/routes"
	"github.com/gin-gonic/gin"
)

func addOpenMiscHandlers(ginEngine *gin.Engine, vars middleware.GinHandlerVars) {
	ginEngine.GET("/",
		routes.Welcome,
	)

	ginEngine.GET("/version",
		func(ctx *gin.Context) {
			routes.VersionInfomation(ctx, vars)
		},
	)

	// ginEngine.GET("/peer-info",
	// 	func(ctx *gin.Context) {
	// 		routes.PeerInfo(ctx, vars)
	// 	},
	// )

	// ginEngine.GET("/hz-info",
	// 	func(ctx *gin.Context) {
	// 		routes.HZInfo(ctx, vars)
	// 	},
	// )

	// ginEngine.GET("/hz-dump",
	// 	func(ctx *gin.Context) {
	// 		routes.HZDump(ctx, vars)
	// 	},
	// )

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

}

func addMiscHandlers(ginEngine *gin.Engine, vars middleware.GinHandlerVars) {
	ginEngine.POST("/checkToken",
		func(ctx *gin.Context) {
			routes.CheckToken(ctx, vars)
		},
	)

	ginEngine.POST("/renewToken",
		func(ctx *gin.Context) {
			routes.RenewToken(ctx, vars)
		},
	)
}
