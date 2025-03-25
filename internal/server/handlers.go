package server

import (
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

var (
	apiKey  string = os.Getenv("API_KEY")
	baseUrl string = os.Getenv("BASE_URL")
)

func isNumber(param string) bool {
	_, err := strconv.Atoi(param)
	return err == nil
}

func sendMalRequest(method string, url string, apiKey string) ([]byte, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Encoding", "gzip")
	req.Header.Add("X-MAL-CLIENT-ID", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, err
	}

	json, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return json, nil
}

func (s *FiberServer) SearchAnimesHandler(c *fiber.Ctx) error {
	q, limit := c.Params("q"), c.Params("limit", "100")
	if q == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing Params")
	}

	url := baseUrl + "/anime?q=" + q + "&limit=" + limit

	json, err := sendMalRequest("GET", url, apiKey)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError)
	}

	return c.Send(json)
}

func (s *FiberServer) GetAnimeDetailsHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing parameters")
	}

	url := baseUrl + "/anime/" + id

	json, err := sendMalRequest("GET", url, apiKey)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError)
	}

	return c.Send(json)
}

func (s *FiberServer) GetSeasonalAnimesHandler(c *fiber.Ctx) error {
	year := c.Params("year")
	season := c.Params("season")
	limit := c.Params("limit", "100")
	allowedSeason := []string{"winter", "spring", "summer", "fall"}

	if year == "" || season == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing parameters")
	} else if !isNumber(year) || !slices.Contains(allowedSeason, season) {
		return fiber.NewError(fiber.StatusBadRequest, "Bad request")
	}

	url := baseUrl + "/anime/season/" + year + "/" + season + "?limit=" + limit

	json, err := sendMalRequest("GET", url, apiKey)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError)
	}

	return c.Send(json)
}
