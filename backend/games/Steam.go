package games

import (
	"TotalControl/backend/steam"
	"TotalControl/backend/utils"
	"net/http"
)

func (g *Game) FetchInfoFromSteam() error {
	// https://store.steampowered.com/api/appdetails?appids=
	if g.ExternalIDs.Steam == "" {
		return nil // No Steam App ID, nothing to fetch
	}

	httpClient := http.DefaultClient

	appInfo, err := steam.GetAppDetails(httpClient, g.ExternalIDs.Steam)
	if err != nil {
		return err
	}

	g.Name = appInfo.Data.Name
	g.Slug = utils.Slugify(appInfo.Data.Name)
	g.Description = appInfo.Data.ShortDescription
	g.Media.Hero = appInfo.Data.CapsuleImage

	return nil
}
