package tests

import (
	"testing"

	"djdees/synthetic_stellar_data/generator"
)

func TestGenerateAllBasic(t *testing.T) {
	cfg := generator.Config{
		NumStars:       10,
		PlanetsPerStar: 5,
		ExoPerStar:     3,
		Seed:           12345,
	}

	data := generator.GenerateAll(cfg)

	if len(data.Stars) != cfg.NumStars {
		t.Errorf("Expected %d stars, got %d", cfg.NumStars, len(data.Stars))
	}

	if len(data.Planets) > cfg.NumStars*cfg.PlanetsPerStar {
		t.Errorf("Too many planets generated: %d", len(data.Planets))
	}

	if len(data.Exoplanets) > cfg.NumStars*cfg.ExoPerStar {
		t.Errorf("Too many exoplanets generated: %d", len(data.Exoplanets))
	}
}

func TestGenerateAllReproducibility(t *testing.T) {
	cfg := generator.Config{
		NumStars:       5,
		PlanetsPerStar: 3,
		ExoPerStar:     2,
		Seed:           99999,
	}

	data1 := generator.GenerateAll(cfg)
	data2 := generator.GenerateAll(cfg)

	if len(data1.Stars) != len(data2.Stars) {
		t.Error("Different number of stars with same seed")
	}

	if len(data1.Planets) != len(data2.Planets) {
		t.Error("Different number of planets with same seed")
	}

	if len(data1.Exoplanets) != len(data2.Exoplanets) {
		t.Error("Different number of exoplanets with same seed")
	}

	// Check first star is identical
	if data1.Stars[0].Name != data2.Stars[0].Name {
		t.Error("Stars not reproducible with same seed")
	}

	if data1.Stars[0].Mass != data2.Stars[0].Mass {
		t.Error("Star mass not reproducible with same seed")
	}
}

func TestStarValidation(t *testing.T) {
	cfg := generator.Config{
		NumStars:       100,
		PlanetsPerStar: 5,
		ExoPerStar:     3,
		Seed:           54321,
	}

	data := generator.GenerateAll(cfg)

	for _, star := range data.Stars {
		// Check UUID format
		if len(star.ID) == 0 {
			t.Error("Star has empty ID")
		}

		// Check positive values
		if star.Mass <= 0 {
			t.Errorf("Star %s has invalid mass: %f", star.Name, star.Mass)
		}

		if star.Radius <= 0 {
			t.Errorf("Star %s has invalid radius: %f", star.Name, star.Radius)
		}

		if star.Temperature <= 0 {
			t.Errorf("Star %s has invalid temperature: %d", star.Name, star.Temperature)
		}

		// Check spectral type format
		if len(star.SpectralType) < 2 {
			t.Errorf("Star %s has invalid spectral type: %s", star.Name, star.SpectralType)
		}
	}
}

func TestPlanetValidation(t *testing.T) {
	cfg := generator.Config{
		NumStars:       50,
		PlanetsPerStar: 8,
		ExoPerStar:     3,
		Seed:           11111,
	}

	data := generator.GenerateAll(cfg)

	for _, planet := range data.Planets {
		// Check UUID
		if len(planet.ID) == 0 {
			t.Error("Planet has empty ID")
		}

		// Check star reference
		if len(planet.StarID) == 0 {
			t.Error("Planet has empty StarID")
		}

		// Check orbital parameters
		if planet.OrbitalPeriod <= 0 {
			t.Errorf("Planet %s has invalid orbital period: %f", planet.Name, planet.OrbitalPeriod)
		}

		if planet.SemiMajorAxis <= 0 {
			t.Errorf("Planet %s has invalid semi-major axis: %f", planet.Name, planet.SemiMajorAxis)
		}

		// Check eccentricity range
		if planet.Eccentricity < 0 || planet.Eccentricity > 1 {
			t.Errorf("Planet %s has invalid eccentricity: %f", planet.Name, planet.Eccentricity)
		}

		// Check physical parameters
		if planet.Mass <= 0 {
			t.Errorf("Planet %s has invalid mass: %f", planet.Name, planet.Mass)
		}

		if planet.Radius <= 0 {
			t.Errorf("Planet %s has invalid radius: %f", planet.Name, planet.Radius)
		}

		// Check discovery year
		if planet.DiscoveryYear < 1990 || planet.DiscoveryYear > 2024 {
			t.Errorf("Planet %s has invalid discovery year: %d", planet.Name, planet.DiscoveryYear)
		}
	}
}

func TestExoplanetValidation(t *testing.T) {
	cfg := generator.Config{
		NumStars:       50,
		PlanetsPerStar: 5,
		ExoPerStar:     5,
		Seed:           22222,
	}

	data := generator.GenerateAll(cfg)

	for _, exo := range data.Exoplanets {
		// Check UUID
		if len(exo.ID) == 0 {
			t.Error("Exoplanet has empty ID")
		}

		// Check star reference
		if len(exo.StarID) == 0 {
			t.Error("Exoplanet has empty StarID")
		}

		// Check detection method
		if len(exo.DetectionMethod) == 0 {
			t.Errorf("Exoplanet %s has empty detection method", exo.Name)
		}

		// Check host distance
		if exo.HostDistance <= 0 {
			t.Errorf("Exoplanet %s has invalid host distance: %f", exo.Name, exo.HostDistance)
		}

		// Check eccentricity
		if exo.Eccentricity < 0 || exo.Eccentricity > 1 {
			t.Errorf("Exoplanet %s has invalid eccentricity: %f", exo.Name, exo.Eccentricity)
		}

		// Check discovery year
		if exo.DiscoveryYear < 1990 || exo.DiscoveryYear > 2024 {
			t.Errorf("Exoplanet %s has invalid discovery year: %d", exo.Name, exo.DiscoveryYear)
		}
	}
}

func TestReferentialIntegrity(t *testing.T) {
	cfg := generator.Config{
		NumStars:       20,
		PlanetsPerStar: 5,
		ExoPerStar:     3,
		Seed:           33333,
	}

	data := generator.GenerateAll(cfg)

	// Create a map of star IDs
	starIDs := make(map[string]bool)
	for _, star := range data.Stars {
		starIDs[star.ID] = true
	}

	// Check all planets reference valid stars
	for _, planet := range data.Planets {
		if !starIDs[planet.StarID] {
			t.Errorf("Planet %s references non-existent star: %s", planet.Name, planet.StarID)
		}
	}

	// Check all exoplanets reference valid stars
	for _, exo := range data.Exoplanets {
		if !starIDs[exo.StarID] {
			t.Errorf("Exoplanet %s references non-existent star: %s", exo.Name, exo.StarID)
		}
	}
}

func BenchmarkGenerateStars(b *testing.B) {
	cfg := generator.Config{
		NumStars:       1000,
		PlanetsPerStar: 8,
		ExoPerStar:     5,
		Seed:           12345,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generator.GenerateAll(cfg)
	}
}
