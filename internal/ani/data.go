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
