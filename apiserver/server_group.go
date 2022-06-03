package apiserver

import (
	"github.com/gin-gonic/gin"
	"github.com/tiki-systems/tikiserver/middleware"
	"github.com/tiki-systems/tikiserver/routes"
)

func addGroupHandlers(ginEngine *gin.Engine, vars middleware.GinHandlerVars) {
	ginEngine.GET("/group",
		func(ctx *gin.Context) {
			routes.ListGroups(ctx, vars)
		},
	)

	ginEngine.GET("/group/:groupName",
		func(ctx *gin.Context) {
			routes.GetGroup(ctx, vars)
		},
	)

	ginEngine.POST("/group/create",
		func(ctx *gin.Context) {
			routes.CreateGroup(ctx, vars)
		},
	)

	ginEngine.DELETE("/group/delete/:groupName",
		func(ctx *gin.Context) {
			routes.DeleteGroup(ctx, vars)
		},
	)

	ginEngine.POST("/group/addmember/:groupName",
		func(ctx *gin.Context) {
			routes.AddMemberToGroup(ctx, vars)
		},
	)

	ginEngine.POST("/group/delmember/:groupName",
		func(ctx *gin.Context) {
			routes.DelMemberFromGroup(ctx, vars)
		},
	)
}
