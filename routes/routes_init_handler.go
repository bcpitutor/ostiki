package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/bcpitutor/ostiki/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var oauthConfGl *oauth2.Config

func GetGoogleAuthConfig() *oauth2.Config {
	return oauthConfGl
}

//func InitHandler(c *gin.Context, logger *zap.Logger, config *appconfig.AppConfig) {
func InitHandler(c *gin.Context, vars middleware.GinHandlerVars) {
	logger := vars.Logger
	config := vars.AppConfig

	logger.Debug("in InitHandler Infomation")

	//redirectUrlBase := viper.GetString("GOOGLE_CB_URL_SERVER")
	// TODO: Confirm w/ Cenk
	redirectUrlBase := config.TikiAuthenticationProviderConfig.RedirectUri
	redirectUrl := redirectUrlBase + "/auth/googlecb"

	oauthConfGl = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  redirectUrl,
		Scopes:       []string{},
		Endpoint:     google.Endpoint,
	}

	//oauthConfGl.ClientID = viper.GetString("GOOGLE_CLIENT_ID")
	oauthConfGl.ClientID = config.TikiAuthenticationProviderConfig.ClientId
	//oauthConfGl.ClientSecret = viper.GetString("GOOGLE_CLIENT_SECRET")
	oauthConfGl.ClientSecret = config.TikiAuthenticationProviderConfig.ClientSecret

	URL, err := url.Parse(oauthConfGl.Endpoint.AuthURL)
	if err != nil {
		logger.Fatal(fmt.Sprintf("ERROR: oAuth config file error: %s", err))
	}

	reqId, _ := uuid.NewRandom()
	port, ok := c.Params.Get("port")
	if !ok {
		port = "3434"
	}

	type Req struct {
		Id   string
		Port string
	}

	var requestInfo = Req{
		Id:   reqId.String(),
		Port: port,
	}

	reqData, err := json.Marshal(requestInfo)
	reqParam := string(reqData)

	parameters := url.Values{}
	parameters.Add("client_id", config.TikiAuthenticationProviderConfig.ClientId)
	scopesByConfig := config.TikiAuthenticationProviderConfig.Scopes
	parameters.Add("scope", strings.Join(scopesByConfig, " "))
	parameters.Add("redirect_uri", oauthConfGl.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("access_type", "offline")
	parameters.Add("prompt", "consent")
	parameters.Add("state", reqParam)

	URL.RawQuery = parameters.Encode()
	url := URL.String()

	logger.Debug(URL.Path)
	logger.Debug(fmt.Sprintf("Setting client_id: %s", oauthConfGl.ClientID))
	logger.Debug(fmt.Sprintf("Setting scope: %s", oauthConfGl.Scopes))
	logger.Debug(fmt.Sprintf("Setting redirect_uri: %s", oauthConfGl.RedirectURL))
	logger.Debug(fmt.Sprintf("Redirecting to: %s", url))

	// Redirect to the google.
	c.Redirect(http.StatusTemporaryRedirect, url)
}
