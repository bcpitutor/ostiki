package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/bcpitutor/ostiki/actions"
	"github.com/bcpitutor/ostiki/appconfig"
	"github.com/bcpitutor/ostiki/repositories"
	"github.com/bcpitutor/ostiki/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/api/idtoken"
)

type GoogleRefreshTokenResponse struct {
	IdToken     string `json:"id_token"`
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type JwtTokenDetails struct {
	At_hash string `json:"at_hash"`
	Aud     string `json:"aud"`
	Azp     string `json:"azp"`
	Email   string `json:"email"`
	Exp     int64  `json:"exp"`
	Hd      string `json:"hd"`
	Iat     int64  `json:"iat"`
	Iss     string `json:"iss"`
	Sub     string `json:"sub"`
}

type GinHandlerVars struct {
	Logger               *zap.Logger
	AWSService           *services.AWSService
	PermissionRepository *repositories.PermissionRepository
	SessionRepository    *repositories.SessionRepository
	BanRepository        *repositories.BanRepository
	DomainRepository     *repositories.DomainRepository
	GroupRepository      *repositories.GroupRepository
	TicketRepository     *repositories.TicketRepository
	AppConfig            *appconfig.AppConfig
}

func Auth(c *gin.Context, vars GinHandlerVars) {
	logger := vars.Logger
	config := vars.AppConfig
	sessionRepository := vars.SessionRepository
	banRepository := vars.BanRepository

	c.Request.Header.Set("email", "")
	c.Request.Header.Set("sub", "")
	c.Request.Header.Set("hd", "")

	reqId, _ := uuid.NewRandom()
	c.Request.Header.Set("x-tsreq-id", reqId.String())

	logger.Sugar().Debugf("Request ID: %s", reqId.String())
	// local development bypass -- TIKISERVER_ENV="local"
	if config.Deployment == "local" {
		logger.Sugar().Info("Local development vars is set")
		if config.DeveloperEmail != "" {
			c.Request.Header.Set("email", config.DeveloperEmail)
		} else {
			c.Request.Header.Set("email", "localdev-tikiserver@itutor.com")
		}

		// if viper.GetString("DEVELOPER_EMAIL") != "" {
		// 	c.Request.Header.Set("email", viper.GetString("DEVELOPER_EMAIL"))
		// } else {
		// 	c.Request.Header.Set("email", "localdev-tikiserver@itutor.com")
		// }
		c.Request.Header.Set("sub", "")
		c.Request.Header.Set("hd", "")
		c.Next()
		return
	}

	// TODO: do some authorization control for user-agent
	authorizationHeader := c.Request.Header.Get("Authorization")

	if authorizationHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"statusCode": http.StatusUnauthorized,
			"message":    "Go and run tikitool auth again.",
			"status":     "error",
			"details":    "-0T",
		})
		return
	}

	logger.Sugar().Debugf("Authorization header: %s", authorizationHeader)

	bearerToken := strings.Split(authorizationHeader, " ")
	if len(bearerToken) < 2 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Go and run tikitool auth again.",
			"status":  "error",
			"details": "",
		})
		return
	}

	// take them under debug..
	// tokenType := bearerToken[0]
	// fmt.Println("arr0:", bearerToken[0])
	// fmt.Println("arr1:", bearerToken[1])

	tokenString := bearerToken[1]

	// // TEST_TOKEN must be set in the environment
	// // Please beware, this skips all the security checks
	// // and never be used in production. It is only for
	// // testing purposes.
	// cfg := appconfig.GetAppConfig()
	// if cfg.Deployment == "local" && tokenString == os.Getenv("TEST_TOKEN") {
	// 	c.Next()
	// }

	//aud := viper.GetString("GOOGLE_CLIENT_ID")
	aud := config.TikiAuthenticationProviderConfig.ClientId

	//hd := viper.GetString("GT_HD")
	hd := config.TikiAuthenticationProviderConfig.GtHd

	// issuers := viper.GetStringSlice("TIKI_GT_ISS")
	// cSecret := viper.GetString("GOOGLE_CLIENT_SECRET")

	// Google auth validation //x
	payload := &idtoken.Payload{
		Issuer:   "",
		Audience: "",
		Claims:   map[string]interface{}{},
	}

	payload, err := idtoken.Validate(context.Background(), tokenString, string(aud))
	if err != nil {
		errMessage := fmt.Sprintf("%s", err)

		logger.Sugar().Debugf("Error validating token: %s", errMessage)

		// Note: this expired token string "idtoken: token expired" comes from google libs.
		if errMessage == "idtoken: token expired" {

			//startRefreshTokenTime := time.Now()
			rToken := c.Request.Header.Get("rtoken")
			logger.Debug("trying to get new rtoken..")
			if rToken == "" {
				logger.Debug("refresh token is not found in the header.")
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "Go and run tikitool auth again.",
					"status":  "error",
					"details": "-5T",
				})
				return
			}

			newToken, err := actions.GetNewGoogleToken(logger, rToken)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "Go and run tikitool auth again.",
					"status":  "error",
					"details": "-1T",
				})
				return
			}

			var gRefTokenResp GoogleRefreshTokenResponse

			byteData, err := json.Marshal(newToken)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "Go and run tikitool auth again.",
					"status":  "error",
					"details": "refresh token marshal error",
				})
				return
			}

			if err := json.Unmarshal(byteData, &gRefTokenResp); err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "Go and run tikitool auth again.",
					"status":  "error",
					"details": "refresh token unmarshal error",
				})
				return
			}

			logger.Sugar().Debugf("expired token renewed and sending... : %s", gRefTokenResp.IdToken)
			payload, err = idtoken.Validate(context.Background(), gRefTokenResp.IdToken, string(aud))
			if err != nil {
				if rToken == "" {
					logger.Sugar().Debugf("Error validating refresh token, rToken is not available: %s", err)
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
						"message": "Go and run tikitool auth again.",
						"status":  "error",
						"details": "-6T",
					})
					return
				}
			}

			// TODO: token is refreshed and new token is received. ctx
			c.Set("newToken", gRefTokenResp.IdToken)
			newTokenData := fmt.Sprintf("Bearer %s", gRefTokenResp.IdToken)
			c.Request.Header.Set("Authorization", newTokenData)
			fmt.Printf("token is refreshed and set in request. %s", newTokenData)

			claims := payload.Claims

			user_email := fmt.Sprintf("%v", claims["email"])
			user_sub := fmt.Sprintf("%v", claims["sub"])
			user_hd := fmt.Sprintf("%v", claims["hd"])

			c.Request.Header.Set("email", user_email)
			c.Request.Header.Set("sub", user_sub)
			c.Request.Header.Set("hd", user_hd)

			// UpdATE session data section -- //
			// TODO: make function for it.
			// bypassed now.
			//controller.UpdateSession(tokenString, gRefTokenResp.IdToken, payload.Expires, rToken)
			sessionRepository.UpdateSession(tokenString, gRefTokenResp.IdToken, payload.Expires, rToken)
			// end of rtoken info

			// end of expired token logic.
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Go and run tikitool auth again.",
				"status":  "error",
				"details": "-1T",
			})
			return
		}

	}

	logger.Sugar().Debugf("Payload : %v", payload)
	if payload.Audience != aud {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Audience mismatch, Go and run tikitool auth again.",
			"status":  "error",
			"details": "-3T",
		})
		return
	}

	// end of google validation control

	// token claims data
	claims := payload.Claims

	logger.Sugar().Debugf("Claims data: %v", claims)
	// TODO:  check against a list of allowed HD's, instead of a fixed one.
	if claims["hd"] != hd {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Domain mismatch, Go and run tikitool auth again.",
			"status":  "error",
			"details": "-4T",
		})
		return
	}

	user_email := fmt.Sprintf("%v", claims["email"])
	user_sub := fmt.Sprintf("%v", claims["sub"])
	user_hd := fmt.Sprintf("%v", claims["hd"])

	if actions.IsUserRevoked(banRepository, logger, user_email) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Your authentication is not valid at the time. Please contact your System Administrator.",
			"status":  "error",
			"details": "-7T",
		})

		logger.Sugar().Debug("Revoked access attempt has been occurred for the user: [%s]", user_email)
		return
	}

	logger.Sugar().Debug("authenticated user email: ", user_email)
	logger.Sugar().Debug("authenticated user hd: ", user_hd)
	logger.Sugar().Debug("authenticated user sub: ", user_sub)

	// putting claims data to the header for controlling some authorization.
	c.Request.Header.Set("email", user_email)
	c.Request.Header.Set("sub", user_sub)
	c.Request.Header.Set("hd", user_hd)

	logger.Sugar().Debugf("Header1 userEmail	: %v\n", c.Request.Header.Get("email"))
	logger.Sugar().Debugf("Header2 userEmail	: %v\n", c.GetHeader("email"))
	logger.Sugar().Debugf("Header userEmail	    : %v\n", c.Request.Header.Get("email"))
	logger.Sugar().Debugf("Header hd 		    : %v \n", c.Request.Header.Get("sub"))
	logger.Sugar().Debugf("Header sub		    : %v \n", c.Request.Header.Get("hd"))

	c.Next()
}
