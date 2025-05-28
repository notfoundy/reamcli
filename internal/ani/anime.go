package ani

import (
	"time"
)

type Anime struct {
	Id                string
	MalId             string
	AnimeDetails      AnimeDetails
	AvailableEpisodes AvailableEpisodes
	Episodes          []Episode
}

type AnimeDetails struct {
	Title          string
	TitleEnglish   string
	Type           string
	Source         string
	Themes         []string
	NumberEpisodes int
	Synopsis       string
	Genres         []string
	Rating         string
	Status         string
	AiredFrom      time.Time
	AiredTo        time.Time
	ImageURL       string
	MalUrl         string
}

type AvailableEpisodes struct {
	Sub int `json:"sub"`
	Dub int `json:"dub"`
	Raw int `json:"raw"`
}

type Episode struct {
	Number  int      `json:"episodeNum"`
	Title   string   `json:"episodeString,omitempty"`
	AirDate string   `json:"airDate,omitempty"`
	Sources []Source `json:"sourceUrls"`
}

type Source struct {
	Name       string  `json:"name"`
	Priority   float64 `json:"priority"`
	StreamType string  `json:"type"`
	EncodedURL string  `json:"sourceUrl"`
}

const (
	allanimeReferer = "https://allmanga.to"
	allanimeBase    = "https://allanime.day"
	allanimeAPI     = "https://api.allanime.day/api"
	userAgent       = "Mozilla/5.0"
)
