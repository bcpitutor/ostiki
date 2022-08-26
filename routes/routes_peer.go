package routes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bcpitutor/ostiki/middleware"
	"github.com/bcpitutor/ostiki/utils"
	"github.com/gin-gonic/gin"
)

func HZDump(c *gin.Context, vars middleware.GinHandlerVars) {
	logger := vars.Logger
	logger.Debug("in HZDump")

	sugar := logger.Sugar()

	imoRepository := vars.ImoRepository
	hzClient := imoRepository.GetHZClient()

	allObj, err := hzClient.GetMap(context.TODO(), "sessions-all")
	if err != nil {
		sugar.Errorf("failed to get map %s: %v", "sessions-all", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("%+v", err),
			"object":  "sessions-all",
		})
		c.Abort()
		return
	}

	activeObj, err := hzClient.GetMap(context.TODO(), "sessions-active")
	if err != nil {
		sugar.Errorf("failed to get map %s: %v", "sessions-active", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("%+v", err),
			"object":  "sessions-active",
		})
		c.Abort()
		return
	}

	expObj, err := hzClient.GetMap(context.TODO(), "sessions-expired")
	if err != nil {
		sugar.Errorf("failed to get map %s: %v", "sessions-expired", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("%+v", err),
			"object":  "sessions-expired",
		})
		c.Abort()
		return
	}

	revokedObj, err := hzClient.GetMap(context.TODO(), "sessions-revoked")
	if err != nil {
		sugar.Errorf("failed to get map %s: %v", "sessions-revoked", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("%+v", err),
			"object":  "sessions-revoked",
		})
		c.Abort()
		return
	}

	allObjES, err := allObj.GetEntrySet(context.TODO())
	if err != nil {
		sugar.Errorf("failed to get entry set %s: %v", "sessions-all", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("%+v", err),
			"object":  "sessions-all",
		})
		c.Abort()
		return
	}

	activeObjES, err := activeObj.GetEntrySet(context.TODO())
	if err != nil {
		sugar.Errorf("failed to get entry set %s: %v", "sessions-active", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("%+v", err),
			"object":  "sessions-active",
		})
		c.Abort()
		return
	}

	expObjES, err := expObj.GetEntrySet(context.TODO())
	if err != nil {
		sugar.Errorf("failed to get entry set %s: %v", "sessions-expired", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("%+v", err),
			"object":  "sessions-expired",
		})
		c.Abort()
		return
	}

	revokedObjES, err := revokedObj.GetEntrySet(context.TODO())
	if err != nil {
		sugar.Errorf("failed to get entry set %s: %v", "sessions-revoked", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("%+v", err),
			"object":  "sessions-revoked",
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"message":      "OK",
		"allObjES":     allObjES,
		"activeObjES":  activeObjES,
		"expObjES":     expObjES,
		"revokedObjES": revokedObjES,
	})

}

func HZInfo(c *gin.Context, vars middleware.GinHandlerVars) {
	logger := vars.Logger
	logger.Debug("in HZ Info")

	imoRepository := vars.ImoRepository
	hzClient := imoRepository.GetHZClient()

	gdoi, err := hzClient.GetDistributedObjectsInfo(context.TODO())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   gdoi,
	})
}

// imo.hzClient.GetDistributedObjectsInfo(context.TODO())
func PeerInfo(c *gin.Context, vars middleware.GinHandlerVars) {
	logger := vars.Logger
	logger.Debug("in Peer Info")

	imoRepository := vars.ImoRepository

	activeObj, err := imoRepository.GetCacheObject("sessions-active")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("%+v", err),
		})
		c.Abort()
		return
	}

	expObj, err := imoRepository.GetCacheObject("sessions-expired")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("%+v", err),
		})
		c.Abort()
		return
	}

	allObj, err := imoRepository.GetCacheObject("sessions-all")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("%+v", err),
		})
		c.Abort()
		return
	}

	myIP := utils.GetOutboundIP().String()

	activeObjES, err := activeObj.GetEntrySet(context.TODO())
	//activeES, err := activeObj.(*hazelcast.Map).GetEntrySet(context.TODO())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("%+v", err),
		})
		c.Abort()
		return
	}

	allObjES, err := allObj.GetEntrySet(context.TODO())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("%+v", err),
		})
		c.Abort()
		return
	}

	// create a string with the entries in activeObjES
	var activeESStr string
	var allObjESStr string
	for _, entry := range activeObjES {
		activeESStr += fmt.Sprintf("%s ", entry.Key)
	}
	for _, entry := range allObjES {
		allObjESStr += fmt.Sprintf("%s ", entry.Key)
	}

	c.JSON(http.StatusOK,
		gin.H{
			"status":            "success",
			"message":           "done",
			"myIP":              myIP,
			"activeSessions":    activeESStr,
			"numActiveSessions": len(activeObjES),
			"allSessions":       allObjESStr,
			"numAllSessions":    len(allObjES),
			"expiredSessions":   expObj,
		},
	)
}
