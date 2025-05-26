package ani

import (
	"github.com/darenliang/jikan-go"
)

func GetSeasonNow() (*jikan.Season, error) {
	season, err := jikan.GetSeasonNow()
	if err != nil {
		return nil, err
	}
	return season, nil
}

// NOTE: Need to check if episodes are updated quickly enough
// Otherwise, I'll need to scrape this from somewhere else
func GetAiredEpisode(id int) (*jikan.AnimeEpisodes, error) {
	episodes, err := jikan.GetAnimeEpisodes(id, 1)
	if err != nil {
		return nil, err
	}
	return episodes, nil
}
