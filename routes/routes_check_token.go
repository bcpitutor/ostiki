package routes

import (
	"encoding/json"
	"net/http"

	"github.com/bcpitutor/ostiki/middleware"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type userInfo struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail string `json:"verifiedEmail"`
	Picture       string `json:"picture"`
	Hd            string `json:"hd"`
}

type RBody struct {
	AccessToken  string    `json:"accessToken"`
	Expire       string    `json:"expire"`
	RefreshToken string    `json:"refreshToken"`
	TokenType    string    `json:"tokenType"`
	IdToken      string    `json:"idToken"`
	UserInfo     *userInfo `json:"userInfo,omitempty"`
}

func CheckToken(c *gin.Context, vars middleware.GinHandlerVars) {
	logger := vars.Logger

	var requestBody RBody
	if err := json.NewDecoder(c.Request.Body).Decode(&requestBody); err != nil {
		logger.Sugar().Errorf("Error decoding request body: %s", err)
	}

	var tokenControl oauth2.Token
	tokenControl.AccessToken = requestBody.AccessToken

	if tokenControl.Valid() {
		// TODO: handle on the client side for keeping things as a shell env or encrypted file etc.
		c.JSON(http.StatusOK, gin.H{
			"message": requestBody,
		})
	} else {
		// TODO: redirect url and client need to be detected...
		c.JSON(http.StatusTemporaryRedirect, gin.H{
			"message": "Redirecting to the login page",
		})
	}
}

func RenewToken(c *gin.Context, vars middleware.GinHandlerVars) {
	logger := vars.Logger
	logger.Info("RenewToken")

	newToken, _ := c.Get("newToken")

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"newToken": newToken,
	})

}
