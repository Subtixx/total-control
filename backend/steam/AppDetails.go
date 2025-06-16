package steam

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func GetAppDetails(httpClient *http.Client, appID string) (*AppDetailsResponse, error) {
	resp, err := httpClient.Get("https://store.steampowered.com/api/appdetails?appids=" + appID)
	if err != nil {
		return nil, errors.New("failed to make request to Steam API")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing response body:", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("received non-200 response from Steam API")
	}

	// The respsonse has a structure like {"<appID>": { "success": true, "data": { ... } }}, we want to ignore the appID key and decode the inner object directly.
	var rawResponse map[string]json.RawMessage
	err = json.NewDecoder(resp.Body).Decode(&rawResponse)
	if err != nil {
		return nil, err
	}
	if len(rawResponse) == 0 {
		return nil, errors.New("empty response from Steam API")
	}
	// Extract the inner data object
	data, ok := rawResponse[appID]
	if !ok {
		return nil, errors.New("app ID not found in response from Steam API")
	}
	// Decode the inner data object into AppDetailsResponse
	var appDetails *AppDetailsResponse
	err = json.Unmarshal(data, &appDetails)
	if err != nil {
		return nil, err
	}
	if appDetails == nil || !appDetails.Success {
		log.Infof("Game is maybe not available in your region or not found: %+v", appDetails)
		return nil, errors.New("app details not found or not available")
	}
	log.Debugf("Fetched app details for app ID %s", appID)

	return appDetails, nil
}
