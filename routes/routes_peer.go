package routes

import (
	"net/http"

	"github.com/bcpitutor/ostiki/middleware"
	"github.com/bcpitutor/ostiki/utils"
	"github.com/gin-gonic/gin"
)

func PeerInfo(c *gin.Context, vars middleware.GinHandlerVars) {
	logger := vars.Logger
	logger.Debug("in Peer Info")

	imoRepository := vars.ImoRepository

	peers := imoRepository.GetPeerIPAddresses()
	myIP := utils.GetOutboundIP().String()

	c.JSON(http.StatusOK,
		gin.H{
			"status":  "success",
			"message": "done",
			"peers":   peers,
			"myIP":    myIP,
		},
	)
}
