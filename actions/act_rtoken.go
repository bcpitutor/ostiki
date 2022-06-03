package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func GetNewGoogleToken(logger *zap.Logger, rToken string) (map[string]any, error) {
	var result map[string]any

	logger.Debug("new token request is starting.")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := http.Client{}
	destUrl := "https://accounts.google.com/o/oauth2/token"

	var URL *url.URL
	URL, err := url.Parse(destUrl)

	params := url.Values{}
	params.Add("client_id", viper.GetString("GOOGLE_CLIENT_ID"))
	params.Add("client_secret", viper.GetString("GOOGLE_CLIENT_SECRET"))
	params.Add("refresh_token", rToken)
	params.Add("grant_type", "refresh_token")

	URL.RawQuery = params.Encode()
	url := URL.String()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		logger.Debug(fmt.Sprintf("new request error: %v", err))
		return result, fmt.Errorf("google token create request error: %s", err)
	}

	res, err := client.Do(req)
	if err != nil {
		logger.Debug(fmt.Sprintf("Google request error: %v", err))
		return result, fmt.Errorf("google token request error: %s", err)
	}

	defer res.Body.Close()

	logger.Sugar().Debugf("Raw Google response: %+v", res)

	rBody, err := io.ReadAll(res.Body)
	if err != nil {
		return result, fmt.Errorf("google token read error: %s", err)
	}

	var tData map[string]any
	if err := json.Unmarshal(rBody, &tData); err != nil {
		return result, fmt.Errorf("google token unmarshall error: %s", err)
	}

	logger.Sugar().Debugf("Google response: %+v", res)

	if tData["id_token"] != nil {
		if tData["id_token"] != "" {
			return tData, nil
		}
	}

	return result, fmt.Errorf("token is not found.")
}
