package ani

import (
	"encoding/json"
	"fmt"
	"time"
)

type RawShow struct {
	Shows struct {
		Edges []struct {
			Id                string `json:"_id"`
			Name              string `json:"name"`
			EnglishName       string `json:"englishName"`
			Thumbnail         string `json:"thumbnail"`
			MalId             string `json:"malId"`
			AvailableEpisodes struct {
				Sub int `json:"sub"`
				Dub int `json:"dub"`
				Raw int `json:"raw"`
			} `json:"availableEpisodes"`
			AvailableEpisodesDetail struct {
				Sub []string `json:"sub"`
				Dub []string `json:"dub"`
				Raw []string `json:"raw"`
			} `json:"availableEpisodesDetail"`
		} `json:"edges"`
	} `json:"shows"`
}

func GetSeasonAnimes(season string, year int) ([]*Anime, error) {
	query := `
  query ($filter: SearchInput!, $page: Int, $limit: Int) {
    shows(search: $filter, page: $page, limit: $limit) {
      edges {
        _id
        name
        englishName
        thumbnail
        malId
        availableEpisodes
        availableEpisodesDetail
      }
    }
  }
`

	vars := map[string]any{
		"filter": map[string]any{
			"season": season,
			"year":   year,
		},
		"page":  1,
		"limit": 20,
	}

	req := Request{
		Query:     query,
		Variables: vars,
	}

	animes := []*Anime{}

	for {
		resp, err := doRequest(allanimeAPI, allanimeReferer, userAgent, req)
		if err != nil {
			return nil, err
		}
		var result RawShow
		if err := json.Unmarshal(resp.Data, &result); err != nil {
			return nil, err
		}

		if len(result.Shows.Edges) == 0 {
			break
		}

		for _, r := range result.Shows.Edges {
			anime := Anime{
				Id:    r.Id,
				MalId: r.MalId,
				AnimeDetails: AnimeDetails{
					Title:          r.Name,
					TitleEnglish:   r.EnglishName,
					Type:           "",
					Source:         "",
					Themes:         []string{},
					NumberEpisodes: 0,
					Synopsis:       "",
					Genres:         []string{},
					Rating:         "",
					Status:         "",
					AiredFrom:      time.Time{},
					AiredTo:        time.Time{},
					ImageURL:       r.Thumbnail,
					MalUrl:         fmt.Sprintf("https://myanimelist.net/anime/%s", r.MalId),
				},
				AvailableEpisodes: AvailableEpisodes{
					Sub: r.AvailableEpisodes.Sub,
					Dub: r.AvailableEpisodes.Dub,
					Raw: r.AvailableEpisodes.Raw,
				},
				Episodes: []*Episode{},
			}
			animes = append(animes, &anime)
		}

		vars["page"] = vars["page"].(int) + 1
	}

	return animes, nil
}
