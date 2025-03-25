package server

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (s *FiberServer) RegisterFiberRoutes() {
	s.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))

	api := s.Group("/api", UnauthMiddleware)

	api.Get("/search/:q/:limit?", s.SearchAnimesHandler)
	api.Get("/anime/:id", s.GetAnimeDetailsHandler)
	api.Get("/season/:year/:season", s.GetSeasonalAnimesHandler)
}
