package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

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
	ImoRepository        *repositories.IMORepository
	AppConfig            *appconfig.AppConfig
	IsCacheReady         *bool
}

func setLocalDeployment(c *gin.Context, config *appconfig.AppConfig, sugar *zap.SugaredLogger) {
	sugar.Infof("Local development vars is set")
	if config.DeveloperEmail != "" {
		c.Request.Header.Set("email", config.DeveloperEmail)
	} else {
		c.Request.Header.Set("email", "localdev-tikiserver@itutor.com")
	}

	c.Request.Header.Set("sub", "")
	c.Request.Header.Set("hd", "")
	c.Next()
	return
}

func parseAuthorizationHeader(authorizationHeader string) (string, error) {
	if authorizationHeader == "" {
		return "", fmt.Errorf("authorization header is empty")
	}
	if !strings.HasPrefix(authorizationHeader, "Baerer ") {
		return "", fmt.Errorf("authorization header is not Baerer")
	}
	return strings.TrimPrefix(authorizationHeader, "Baerer "), nil
}

func Auth(c *gin.Context, vars GinHandlerVars) {
	logger := vars.Logger
	config := vars.AppConfig
	sugar := logger.Sugar() // TODO: use for all logging instead of logger
	sessionRepository := vars.SessionRepository
	banRepository := vars.BanRepository
	imoRepository := vars.ImoRepository

	//sugar.Infof("Auth middleware called")
	c.Request.Header.Set("email", "")
	c.Request.Header.Set("sub", "")
	c.Request.Header.Set("hd", "")

	// a := c.Request.Header
	// keys := make([]string, 0, len(a))
	// for k := range a {
	// 	keys = append(keys, k)
	// }
	// //sugar.Infof("Request headers: %+v", keys)
	// for _, k := range keys {
	// 	sugar.Infof("Key: %s, Value: %s", k, a.Get(k))
	// }

	reqId, _ := uuid.NewRandom()
	c.Request.Header.Set("x-tsreq-id", reqId.String())

	//sugar.Infof("Request ID: %s", reqId.String())
	// local development bypass -- TIKISERVER_ENV="local"
	if config.Deployment == "local" {
		setLocalDeployment(c, config, sugar)
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

	//logger.Sugar().Debugf("Authorization header: %s", authorizationHeader)
	//bearerToken := strings.Split(authorizationHeader, " ")
	tokenString, err := parseAuthorizationHeader(authorizationHeader)
	if err != nil {
		sugar.Debugf("Error parsing authorization header: %s", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Go and run tikitool auth again.",
			"status":  "error",
			"details": err,
		})
		return
	}
	//sugar.Debugf("Got Token string: %s", tokenString)

	sessionByToken, err := sessionRepository.GetSessionByToken(tokenString)
	//sessionByToken, err := imoRepository.GetSessionByToken(tokenString)
	if err != nil {
		logger.Sugar().Debugf("Error getting session by token: %s", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Your session couldn't been found. Please login again.",
			"details": fmt.Sprintf("%v", err),
		})
		return
	}
	//logger.Sugar().Debugf("sessionsByToken: %+v", sessionByToken)

	// TODO: DO NOT FORGET TO REMOVE !=0 SINCE IT'S ONLY THERE FOR COMPABILITY WITH OLD SERVER
	if sessionByToken.SessionExpEpoch != 0 && sessionByToken.SessionExpEpoch < time.Now().Unix() {
		logger.Sugar().Debug("Session %s has expired.", sessionByToken.SessID)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Your session has expired. Please login again.",
			"details": fmt.Sprintf("%v", err),
		})
		return
	}
	//logger.Sugar().Debugf("Session %s is valid.", sessionByToken.SessID)

	// take them under debug..
	// tokenType := bearerToken[0]
	// fmt.Println("arr0:", bearerToken[0])
	// fmt.Println("arr1:", bearerToken[1])

	//tokenString := bearerToken[1]

	// // TEST_TOKEN must be set in the environment
	// // Please beware, this skips all the security checks
	// // and never be used in production. It is only for
	// // testing purposes.
	// cfg := appconfig.GetAppConfig()
	// if cfg.Deployment == "local" && tokenString == os.Getenv("TEST_TOKEN") {
	// 	c.Next()
	// }

	aud := config.TikiAuthenticationProviderConfig.ClientId
	hd := config.TikiAuthenticationProviderConfig.GtHd

	// issuers := viper.GetStringSlice("TIKI_GT_ISS")
	// cSecret := viper.GetString("GOOGLE_CLIENT_SECRET")

	// Google auth validation //x
	payload := &idtoken.Payload{
		Issuer:   "",
		Audience: "",
		Claims:   map[string]interface{}{},
	}

	// TODO: Remove
	if c.Request.Header.Get("Hede") == "hodo" {
		sugar.Infof("Token renewal requested")
		_, err := attemptToRenewToken(config, sugar, sessionRepository, imoRepository, c, aud, tokenString)
		if err != nil {
			sugar.Debugf("Error renewing token: %s", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Go and run tikitool auth again.",
				"status":  "error",
				"details": err,
			})
			return
		}
		sugar.Infof("Token renewed")
		c.Next()
		return
	}

	payload, err = idtoken.Validate(context.Background(), tokenString, string(aud))
	if err != nil {
		errMessage := fmt.Sprintf("%s", err)
		logger.Sugar().Debugf("Error validating token: %s", errMessage)

		// Note: this expired token string "idtoken: token expired" comes from google libs.
		if errMessage == "idtoken: token expired" {
			sugar.Debugf("Token is expired, attempting to renew")
			p, err := attemptToRenewToken(config, sugar, sessionRepository, imoRepository, c, aud, tokenString)
			if err != nil {
				sugar.Debugf("Error renewing token: %s", err)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "Go and run tikitool auth again.",
					"status":  "error",
					"details": err,
				})
				return
			}
			payload = p
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Go and run tikitool auth again.",
				"status":  "error",
				"details": "-1T",
			})
			return
		}

	}

	//logger.Sugar().Debugf("Payload : %+v", payload)
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

	//logger.Sugar().Debugf("Claims data: %v", claims)
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

	// logger.Sugar().Debug("authenticated user email: ", user_email)
	// logger.Sugar().Debug("authenticated user hd: ", user_hd)
	// logger.Sugar().Debug("authenticated user sub: ", user_sub)

	// tokenString
	// sessionsByToken := imoRepository.GetSessionsByToken(idToken.(string))
	// // detect if sessions expired
	// if sessionsByToken[0].SessionExpEpoch < time.Now().Unix() {
	// 	sugar.Debugf("Found expired session for user [%s] ID: [%s]", sessionsByToken[0].SessionOwner, sessionsByToken[0].SessID)
	// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
	// 		"message": "Your session has expired. Please login again.",
	// 		"details": fmt.Sprintf("%v", err),
	// 	})
	// 	return
	// }

	// putting claims data to the header for controlling some authorization.
	c.Request.Header.Set("email", user_email)
	c.Request.Header.Set("sub", user_sub)
	c.Request.Header.Set("hd", user_hd)

	// logger.Sugar().Debugf("Header1 userEmail: %v\n", c.Request.Header.Get("email"))
	// logger.Sugar().Debugf("Header2 userEmail: %v\n", c.GetHeader("email"))
	// logger.Sugar().Debugf("Header userEmail : %v\n", c.Request.Header.Get("email"))
	// logger.Sugar().Debugf("Header hd 		: %v\n", c.Request.Header.Get("sub"))
	// logger.Sugar().Debugf("Header sub		: %v\n", c.Request.Header.Get("hd"))

	sugar.Infof("User [%s] authenticated successfully. ReqId: %s", user_email, reqId.String())
	c.Next()
}

func attemptToRenewToken(config *appconfig.AppConfig, sugar *zap.SugaredLogger, sessionRepository *repositories.SessionRepository, imoRepository *repositories.IMORepository, c *gin.Context, aud string, tokenString string) (*idtoken.Payload, error) {
	sugar.Infof("attemptToRenewToken")
	rToken := c.Request.Header.Get("rtoken")
	sugar.Debug("trying to get new rtoken..")
	if rToken == "" {
		// sugar.Debug("refresh token is not found in the header.")
		// c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		// 	"message": "Go and run tikitool auth again.",
		// 	"status":  "error",
		// 	"details": "-5T",
		// })
		return nil, fmt.Errorf("refresh token is not found in the header.")
	}

	sugar.Debugf("calling google to get new token with refresh token")
	newToken, err := actions.GetNewGoogleToken(config, sugar, rToken)
	if err != nil {
		// sugar.Debugf("Error getting new token: %s", err)
		// c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		// 	"message": "Go and run tikitool auth again.",
		// 	"status":  "error",
		// 	"details": "-1T",
		// })
		return nil, fmt.Errorf("Error getting new token: %s", err)
	}
	//sugar.Debug("new token is generated: ", newToken)

	var gRefTokenResp GoogleRefreshTokenResponse

	byteData, err := json.Marshal(newToken)
	if err != nil {
		// c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		// 	"message": "Go and run tikitool auth again.",
		// 	"status":  "error",
		// 	"details": "refresh token marshal error",
		// })
		return nil, fmt.Errorf("Error marshal new token: %s", err)
	}

	if err := json.Unmarshal(byteData, &gRefTokenResp); err != nil {
		// c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		// 	"message": "Go and run tikitool auth again.",
		// 	"status":  "error",
		// 	"details": "refresh token unmarshal error",
		// })
		return nil, fmt.Errorf("Error unmarshal new token: %s", err)
	}

	//sugar.Debugf("expired token renewed and sending... : %s", gRefTokenResp.IdToken)
	payload, err := idtoken.Validate(context.Background(), gRefTokenResp.IdToken, string(aud))
	if err != nil {
		if rToken == "" {
			sugar.Debugf("Error validating refresh token, rToken is not available: %s", err)
			// c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			// 	"message": "Go and run tikitool auth again.",
			// 	"status":  "error",
			// 	"details": "-6T",
			// })
			return nil, fmt.Errorf("Error validating refresh token, rToken is not available: %s", err)
		}
	}

	// TODO: token is refreshed and new token is received. ctx
	c.Set("newToken", gRefTokenResp.IdToken)
	newTokenData := fmt.Sprintf("Baerer %s", gRefTokenResp.IdToken)
	c.Request.Header.Set("Authorization", newTokenData)
	//fmt.Printf("token is refreshed and set in request. %s", newTokenData)

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

	// TODO: HANDLE ERROR HERE
	sessionRepository.UpdateSession(tokenString, gRefTokenResp.IdToken, payload.Expires, rToken)
	//sessions := imoRepository.

	// end of rtoken info

	// end of expired token logic.

	return payload, nil
}
