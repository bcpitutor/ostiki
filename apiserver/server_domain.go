package apiserver

import (
	"github.com/gin-gonic/gin"
	"github.com/tiki-systems/tikiserver/middleware"
	"github.com/tiki-systems/tikiserver/routes"
)

func addDomainHandlers(ginEngine *gin.Engine, vars middleware.GinHandlerVars) {
	ginEngine.GET("/domain/list", func(ctx *gin.Context) {
		routes.ListDomains(ctx, vars)
	})

	ginEngine.POST("/domain/create", func(ctx *gin.Context) {
		routes.CreateDomain(ctx, vars)
	})

	ginEngine.POST("/domain/get", func(ctx *gin.Context) {
		routes.GetDomain(ctx, vars)
	})

	ginEngine.DELETE("/domain/delete/:domainPath", func(ctx *gin.Context) {
		routes.DeleteDomain(ctx, vars)
	})
}
