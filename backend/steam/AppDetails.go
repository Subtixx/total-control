package steam

import (
	"TotalControl/backend/utils"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path"
)

var appDetailsCache *utils.Cache

func InitializeAppDetailsCache() {
	if appDetailsCache != nil {
		log.Debug("App details cache already initialized")
		return
	}

	cachePath := utils.GetAppCachePath()
	steamCachePath := path.Join(cachePath, "steam.json")
	if _, err := os.Stat(steamCachePath); os.IsNotExist(err) {
		if err := utils.FileTouch(steamCachePath, "{}"); err != nil {
			log.Errorf("Failed to create cache file %s: %v", steamCachePath, err)
			return
		}
	}
	appDetailsCache = utils.NewCache(steamCachePath, "steam_app_details")
}

func GetAppDetails(httpClient *http.Client, appID string) (*AppDetailsResponse, error) {
	if appDetailsCache == nil {
		InitializeAppDetailsCache()
		if appDetailsCache == nil {
			return nil, errors.New("app details cache is not initialized")
		}
	}

	// Cache
	cachedAppDetail, err := appDetailsCache.Get(appID)
	if err == nil && cachedAppDetail != nil {
		log.Debugf("Cache hit for app ID %s", appID)
		return convertToAppDetailsResponse(cachedAppDetail), nil
	}
	log.Debugf("Cache miss for app ID %s, fetching from Steam API", appID)

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

	// Cache the app details
	err = appDetailsCache.Set(appID, appDetails, utils.DefaultCacheTTL)
	if err != nil {
		log.Errorf("Failed to cache app details for app ID %s: %v", appID, err)
	}

	err = appDetailsCache.Save()
	if err != nil {
		log.Errorf("Failed to save app details cache for app ID %s: %v", appID, err)
	}

	return appDetails, nil
}

// convertToAppDetailsResponse attempts to convert a cached value to *AppDetailsResponse
func convertToAppDetailsResponse(data interface{}) *AppDetailsResponse {
	if resp, ok := data.(*AppDetailsResponse); ok {
		return resp
	}
	// If it's a map, marshal and unmarshal to the correct type
	if m, ok := data.(map[string]interface{}); ok {
		b, err := json.Marshal(m)
		if err == nil {
			var appDetails AppDetailsResponse
			if err := json.Unmarshal(b, &appDetails); err == nil {
				return &appDetails
			}
		}
	}
	return nil
}
