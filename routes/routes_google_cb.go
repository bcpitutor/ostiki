package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/bcpitutor/ostiki/middleware"
	"github.com/bcpitutor/ostiki/models"
	"github.com/gin-gonic/gin"
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
	imoRepository := vars.ImoRepository
	config := vars.AppConfig
	sugar := logger.Sugar()

	sugar.Debug("in GoogleCBHandler")

	// TODO: Handle if c.Query has variables (to prevent nil pointer error)
	code := c.Query("code")
	if code == "" {
		sugar.Errorf("Got empty code from Google: [%+v]", c)
	}

	scope := c.Query("scope")
	authUser := c.Query("authuser")
	hostDomain := c.Query("hd")
	// TODO: check this state if it matches what we expect to detect unauthorized clients
	state := c.Query("state")

	// TODO: Establish flags mechanism for logging levels.
	sugar.Debugf("Got code from Google: %s", code)
	sugar.Debugf("Scope: %s", scope)
	sugar.Debugf("Auth User: %s", authUser)
	sugar.Debugf("HD: %s", hostDomain)
	sugar.Debugf("State: %s", state)

	tokenResponse, err := GetGoogleAuthConfig().Exchange(context.TODO(), code)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Request failed try again. Contact your system administrator.",
			"details": fmt.Sprintf("%v", err),
		})
		return
	}

	sugar.Debugf("full response of exchange: |%+v|", tokenResponse)

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
	sugar.Infof("UserInfo : %+v", respBody)

	// get users sessions from imo
	sessionsByEmail := imoRepository.GetSessionsByEmail(respBody.Email)
	sugar.Infof("Got %d sessions for user: %s", len(sessionsByEmail), respBody.Email)

	if len(sessionsByEmail) > config.SessionMaxSimultaneousUsers {
		// delete users all sessions except the 3 newest ones
		sugar.Infof("Deleting %d sessions for user: %s", len(sessionsByEmail)-config.SessionMaxSimultaneousUsers, respBody.Email)
		for _, session := range sessionsByEmail[config.SessionMaxSimultaneousUsers:] {
			// first delete from the db instead of imo, then if the delete successful delete from imo
			err := sessionRepository.DeleteSession(session.SessID, session.Epoch)
			if err != nil {
				sugar.Errorf("Failed to delete session from DB: %s", session.SessID)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "Request failed try again. Contact your system administrator. Error: " + err.Error(),
					"details": fmt.Sprintf("%v", err),
				})
				return
			}

			imoRepository.DeleteSession(session.SessID, false)
		}
		peerMsg := "upd, session"
		imoRepository.SendMessageToPeers(peerMsg, imoRepository.PeerIPAddresses)
	}

	// detect if sessions expired
	for _, session := range sessionsByEmail {
		if session.SessionExpEpoch < time.Now().Unix() {
			sugar.Debugf("Found expired session for user [%s] ID: [%s]", session.SessionOwner, session.SessID)
		}
	}

	type Res struct {
		Id   string
		Port string
	}
	var stateResp Res

	if err := json.Unmarshal([]byte(state), &stateResp); err != nil {
		sugar.Infof("State decoding error: %s", state)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Request failed, state decoding unmarshall error.",
			"details": fmt.Sprintf("%v", err),
		})
		return
	}

	sugar.Info("state response info:", stateResp)

	// Add to the db.
	var sessInfo models.Session

	now := time.Now().Unix()

	sessInfo.SessID = stateResp.Id
	sessInfo.SessionOwner = respBody.Email
	sessInfo.IdToken = fmt.Sprintf("%v", idToken)
	sessInfo.AccessToken = tokenResponse.AccessToken
	sessInfo.Expire = strconv.FormatInt(tokenResponse.Expiry.Unix(), 10)
	sessInfo.RefreshToken = tokenResponse.RefreshToken
	sessInfo.TokenType = tokenResponse.TokenType
	sessInfo.UserInfo = respBody
	sessInfo.Epoch = now                                     // time that the session was created
	sessInfo.SessionExpEpoch = now + config.SessionMaxLength // time that the session will expire

	err = sessionRepository.CreateSession(&sessInfo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "request failed try again.",
			"details": fmt.Sprintf("%v", err),
		})
		return
	}

	sugar.Debugf("Session created: %+v", sessInfo)
	err = imoRepository.AddSession(sessInfo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "request failed try again.",
			"details": fmt.Sprintf("%v", err),
		})
		return
	}

	// activeSessions := imoRepository.GetSessions()
	// imoRepository.SetSessions(append(activeSessions, sessInfo))

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
