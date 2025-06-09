package steam

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"strconv"
)

func GetAppDetails(appID int) (*AppDetailsResponse, error) {
	resp, err := http.Get("https://store.steampowered.com/api/appdetails?appids=" + strconv.Itoa(appID))
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
	data, ok := rawResponse[strconv.Itoa(appID)]
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
	log.Debugf("Fetched app details for app ID %d", appID)

	// Save rawResponse to a file for debugging purposes
	file, err := os.Create(fmt.Sprintf("data/steam/app_details_%d.json", appID))
	if err != nil {
		log.Errorf("Failed to create file for app details: %v", err)
	} else {
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Errorf("Error closing file for app details: %v", err)
			}
		}(file)

		// Write the raw response to the file
		if _, err := file.Write(data); err != nil {
			log.Errorf("Failed to write app details to file: %v", err)
		} else {
			log.Infof("App details saved to data/steam/app_details_%d.json", appID)
		}
	}

	return appDetails, nil
}
