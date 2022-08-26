package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/bcpitutor/ostiki/appconfig"
	"go.uber.org/zap"
)

func GetNewGoogleToken(config *appconfig.AppConfig, sugar *zap.SugaredLogger, rToken string) (map[string]any, error) {
	var result map[string]any

	sugar.Debug("new token request is starting.")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := http.Client{}
	destUrl := "https://accounts.google.com/o/oauth2/token"

	var URL *url.URL
	URL, err := url.Parse(destUrl)

	params := url.Values{}
	//params.Add("client_id", viper.GetString("GOOGLE_CLIENT_ID"))
	params.Add("client_id", config.TikiAuthenticationProviderConfig.ClientId)
	//params.Add("client_secret", viper.GetString("GOOGLE_CLIENT_SECRET"))
	params.Add("client_secret", config.TikiAuthenticationProviderConfig.ClientSecret)
	params.Add("refresh_token", rToken)
	params.Add("grant_type", "refresh_token")

	URL.RawQuery = params.Encode()
	url := URL.String()

	//sugar.Debugf("URL: %s", url)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		sugar.Debug(fmt.Sprintf("new request error: %v", err))
		return result, fmt.Errorf("google token create request error: %s", err)
	}

	res, err := client.Do(req)
	if err != nil {
		sugar.Debug(fmt.Sprintf("Google request error: %v", err))
		return result, fmt.Errorf("google token request error: %s", err)
	}

	defer res.Body.Close()

	//sugar.Debugf("Raw Google response: %+v", res)

	rBody, err := io.ReadAll(res.Body)
	if err != nil {
		return result, fmt.Errorf("google token read error: %s", err)
	}

	var tData map[string]any
	if err := json.Unmarshal(rBody, &tData); err != nil {
		sugar.Debugf("google token unmarshal error: %v", err)
		return result, fmt.Errorf("google token unmarshall error: %s", err)
	}

	//sugar.Debugf("Google response: %+v", res)
	//sugar.Debugf("TData: %+v", tData)

	if tData["error"] != nil {
		sugar.Debugf("Google token error: %+v", tData["error"])
		return result, fmt.Errorf("%s", tData["error_description"])
	}
	//sugar.Debugf("No error, Google token response: %+v", tData)

	return tData, nil
	// if tData["id_token"] != nil {
	// 	sugar.Debugf("1 Google response: %+v", tData["id_token"])
	// 	if tData["id_token"] != "" {
	// 		sugar.Debugf("2 Google response: %+v", tData["id_token"])
	// 		return tData, nil
	// 	}
	// }

	// return result, fmt.Errorf("token is not found.")
}
