package apiserver

import (
	"github.com/bcpitutor/ostiki/middleware"
	"github.com/bcpitutor/ostiki/routes"
	"github.com/gin-gonic/gin"
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
