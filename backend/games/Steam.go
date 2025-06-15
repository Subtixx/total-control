package games

import (
	"TotalControl/backend/steam"
	"TotalControl/backend/utils"
)

func (g *Game) FetchInfoFromSteam() error {
	// https://store.steampowered.com/api/appdetails?appids=
	if g.ExternalIDs.Steam == "" {
		return nil // No Steam App ID, nothing to fetch
	}

	appInfo, err := steam.GetAppDetails(g.ExternalIDs.Steam)
	if err != nil {
		return err
	}

	g.Name = appInfo.Data.Name
	g.Slug = utils.Slugify(appInfo.Data.Name)
	g.Description = appInfo.Data.ShortDescription
	g.Media.Hero = appInfo.Data.CapsuleImage

	return nil
}
