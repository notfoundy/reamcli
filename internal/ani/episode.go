package ani

import (
	"encoding/json"
	"strconv"
	"time"
)

type RawEpisode struct {
	Episode struct {
		EpisodeString string `json:"episodeString"`
		SourceUrls    []struct {
			ClassName string `json:"className"`
			Downloads *struct {
				DownloadUrl string `json:"downloadUrl"`
				SourceName  string `json:"sourceName"`
			} `json:"downloads,omitempty"`
			Priority   float64 `json:"priority"`
			SourceName string  `json:"sourceName"`
			SourceUrl  string  `json:"sourceUrl"`
			StreamerId string  `json:"streamerId"`
			Type       string  `json:"type"`
			Sandbox    string  `json:"sandbox,omitempty"`
		} `json:"sourceUrls"`
		UploadDate struct {
			Year   int `json:"year"`
			Month  int `json:"month"`
			Date   int `json:"date"`
			Hour   int `json:"hour"`
			Minute int `json:"minute"`
			Second int `json:"second"`
		} `json:"uploadDate"`
	} `json:"episode"`
}

func GetEpisodesAnimes(anime *Anime, translation string) ([]*Episode, error) {
	query := `
	query ($showId: String!, $translationType: VaildTranslationTypeEnumType!, $episodeString: String!) {
		episode(showId: $showId, translationType: $translationType, episodeString: $episodeString) {
			episodeString
			sourceUrls
			uploadDate
		}
	}
	`

	vars := map[string]any{
		"showId":          anime.Id,
		"translationType": translation,
		"episodeString":   "1",
	}

	req := Request{
		Query:     query,
		Variables: vars,
	}

	for i := range anime.AvailableEpisodes.Sub {
		resp, err := doRequest(allanimeAPI, allanimeReferer, userAgent, req)
		if err != nil {
			return nil, err
		}

		var result RawEpisode
		if err := json.Unmarshal(resp.Data, &result); err != nil {
			return nil, err
		}

		sources := []*Source{}
		for _, s := range result.Episode.SourceUrls {
			var url, name string
			if s.SourceUrl != "" {
				url = s.SourceUrl
				name = s.SourceName
			} else {
				url = s.Downloads.DownloadUrl
				name = s.Downloads.SourceName
			}
			source := Source{
				Name:      name,
				Priority:  s.Priority,
				Type:      s.Type,
				SourceUrl: url,
			}
			sources = append(sources, &source)
		}

		date := result.Episode.UploadDate
		episode := &Episode{
			Number:    i,
			AiredDate: time.Date(date.Year, time.Month(date.Month), date.Date, 0, 0, 0, 0, time.UTC),
			Sources:   sources,
		}

		anime.Episodes = append(anime.Episodes, episode)
		vars["episodeString"] = strconv.Itoa(i)
	}

	episodes := anime.Episodes
	return episodes, nil
}
