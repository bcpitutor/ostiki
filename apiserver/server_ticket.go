package apiserver

import (
	"github.com/gin-gonic/gin"
	"github.com/tiki-systems/tikiserver/middleware"
	"github.com/tiki-systems/tikiserver/routes"
)

func addTicketHandlers(ginEngine *gin.Engine, vars middleware.GinHandlerVars) {
	ginEngine.GET("/ticket/list",
		func(ctx *gin.Context) {
			routes.ListTickets(ctx, vars)
		},
	)

	ginEngine.POST("/ticket/create",
		func(ctx *gin.Context) {
			routes.CreateTicket(ctx, vars)
		},
	)

	ginEngine.POST("/ticket/get",
		func(ctx *gin.Context) {
			routes.GetTicket(ctx, vars)
		},
	)

	ginEngine.DELETE("/ticket/delete",
		func(ctx *gin.Context) {
			routes.DeleteTicket(ctx, vars)
		},
	)

	ginEngine.POST("/ticket/obtain",
		func(ctx *gin.Context) {
			routes.ObtainTicket(ctx, vars)
		},
	)

	ginEngine.POST("/ticket/secret/set",
		func(ctx *gin.Context) {
			routes.TicketSetSecret(ctx, vars)
		},
	)

	ginEngine.POST("/ticket/secret/get",
		func(ctx *gin.Context) {
			routes.TicketGetSecret(ctx, vars)
		},
	)
}
