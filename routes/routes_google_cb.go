package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tiki-systems/tikiserver/middleware"
	"github.com/tiki-systems/tikiserver/models"
)

type UserInfoRespBody struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail string `json:"verifiedEmail"`
	Picture       string `json:"picture"`
	Hd            string `json:"hd"`
}

func GoogleCBHandler(c *gin.Context, vars middleware.GinHandlerVars) {
	logger := vars.Logger
	sessionRepository := vars.SessionRepository

	//models.GoogleCBVars
	logger.Debug("in InitHandler Infomation")

	// TODO: Handle if c.Query has variables (to prevent nil pointer error)
	code := c.Query("code")
	if code == "" {
		logger.Error(fmt.Sprintf("Got empty code from Google: [%+v]", c))
	}

	scope := c.Query("scope")
	authUser := c.Query("authuser")
	hostDomain := c.Query("hd")
	// TODO: check this state if it matches what we expect to detect unauthorized clients
	state := c.Query("state")

	// TODO: Establish flags mechanism for logging levels.
	logger.Debug(fmt.Sprintf("Got code from Google: %s", code))
	logger.Debug(fmt.Sprintf("Scope: %s", scope))
	logger.Debug(fmt.Sprintf("Auth User: %s", authUser))
	logger.Debug(fmt.Sprintf("HD: %s", hostDomain))
	logger.Debug(fmt.Sprintf("State: %s", state))

	tokenResponse, err := GetGoogleAuthConfig().Exchange(context.TODO(), code)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Request failed try again. Contact your system administrator.",
			"details": fmt.Sprintf("%v", err),
		})
		return
	}

	logger.Sugar().Debugf("full response of exchange: |%+v|", tokenResponse)

	idToken := tokenResponse.Extra("id_token")

	var httpClient = &http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(tokenResponse.AccessToken))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "request failed try again.",
			"details": fmt.Sprintf("%v", err),
		})
		return
	}
	var respBody UserInfoRespBody
	json.NewDecoder(resp.Body).Decode(&respBody)

	defer resp.Body.Close()
	logger.Info(fmt.Sprintf("UserInfo : %+v", respBody))

	type Res struct {
		Id   string
		Port string
	}
	var stateResp Res

	if err := json.Unmarshal([]byte(state), &stateResp); err != nil {
		logger.Info(fmt.Sprintf("State decoding error: %s", state))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Request failed, state decoding unmarshall error.",
			"details": fmt.Sprintf("%v", err),
		})
		return
	}

	logger.Sugar().Info("state response info:", stateResp)

	// Add to the db.
	var sessInfo models.Session

	sessInfo.SessID = stateResp.Id
	sessInfo.SessionOwner = respBody.Email
	sessInfo.IdToken = fmt.Sprintf("%v", idToken)
	sessInfo.AccessToken = tokenResponse.AccessToken
	sessInfo.Expire = strconv.FormatInt(tokenResponse.Expiry.Unix(), 10)
	sessInfo.RefreshToken = tokenResponse.RefreshToken
	sessInfo.TokenType = tokenResponse.TokenType
	sessInfo.UserInfo = respBody
	sessInfo.Epoch = time.Now().Unix() + (1 * 12 * 3600)

	err = sessionRepository.CreateSession(&sessInfo)
	//dResp, err := controller.CreateSess(sessInfo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "request failed try again.",
			"details": fmt.Sprintf("%v", err),
		})
		return
	}

	logger.Sugar().Debugf("Session created: %+v", sessInfo)
	c.HTML(
		http.StatusOK,
		"index.html",
		gin.H{
			"token":        sessInfo.IdToken,
			"rToken":       sessInfo.RefreshToken,
			"clientID":     GetGoogleAuthConfig().ClientID,
			"clientSecret": GetGoogleAuthConfig().ClientSecret,
			"email":        respBody.Email,
			"port":         stateResp.Port,
		},
	)
}
