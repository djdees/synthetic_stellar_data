package generator

import (
	"djdees/synthetic_stellar_data/models"
)

// Config holds the configuration for data generation
type Config struct {
	NumStars       int   // Number of stars to generate
	PlanetsPerStar int   // Maximum number of planets per star
	ExoPerStar     int   // Maximum number of exoplanets per star
	Seed           int64 // Random seed for reproducibility
}

// GeneratedData holds all generated entities
type GeneratedData struct {
	Stars      []models.Star
	Planets    []models.Planet
	Exoplanets []models.Exoplanet
}
