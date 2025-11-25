package generator

import (
	"fmt"
	"math"
	"math/rand"

	"djdees/synthetic_stellar_data/models"
	"github.com/google/uuid"
)

// Spectral types with their characteristics
var spectralTypes = []struct {
	class       string
	massRange   [2]float64 // Min, Max in solar masses
	radiusRange [2]float64 // Min, Max in solar radii
	tempRange   [2]int32   // Min, Max in Kelvin
}{
	{"O", [2]float64{16.0, 90.0}, [2]float64{6.6, 15.0}, [2]int32{30000, 50000}},
	{"B", [2]float64{2.1, 16.0}, [2]float64{1.8, 6.6}, [2]int32{10000, 30000}},
	{"A", [2]float64{1.4, 2.1}, [2]float64{1.4, 1.8}, [2]int32{7500, 10000}},
	{"F", [2]float64{1.04, 1.4}, [2]float64{1.15, 1.4}, [2]int32{6000, 7500}},
	{"G", [2]float64{0.8, 1.04}, [2]float64{0.96, 1.15}, [2]int32{5200, 6000}},
	{"K", [2]float64{0.45, 0.8}, [2]float64{0.7, 0.96}, [2]int32{3700, 5200}},
	{"M", [2]float64{0.08, 0.45}, [2]float64{0.1, 0.7}, [2]int32{2400, 3700}},
}

var atmosphereTypes = []string{
	"H2/He dominant",
	"N2/O2 dominant",
	"CO2 dominant",
	"Thin atmosphere",
	"No atmosphere",
	"H2/He with methane",
	"Sulfuric compounds",
	"Water vapor rich",
}

var detectionMethods = []string{
	"Transit",
	"Radial Velocity",
	"Direct Imaging",
	"Gravitational Microlensing",
	"Astrometry",
	"Transit Timing Variation",
}

// randFloat generates a random float64 between min and max
func randFloat(r *rand.Rand, min, max float64) float64 {
	return min + r.Float64()*(max-min)
}

// randInt generates a random int32 between min and max
func randInt(r *rand.Rand, min, max int32) int32 {
	return min + r.Int31n(max-min+1)
}

// generateSpectralType generates a realistic spectral type
func generateSpectralType(r *rand.Rand) string {
	// Weight spectral types by frequency (M stars are most common)
	weights := []float64{0.00003, 0.13, 2.0, 3.0, 7.6, 12.1, 76.45}
	total := 0.0
	for _, w := range weights {
		total += w
	}

	val := r.Float64() * total
	cumulative := 0.0
	classIdx := 0

	for i, w := range weights {
		cumulative += w
		if val <= cumulative {
			classIdx = i
			break
		}
	}

	class := spectralTypes[classIdx].class
	subclass := r.Intn(10)
	luminosity := "V" // Main sequence

	return fmt.Sprintf("%s%d%s", class, subclass, luminosity)
}

// generateStar creates a realistic star
func generateStar(r *rand.Rand, index int) models.Star {
	spectralType := generateSpectralType(r)
	classChar := spectralType[0]

	// Find spectral class characteristics
	var classIdx int
	for i, st := range spectralTypes {
		if st.class[0] == classChar {
			classIdx = i
			break
		}
	}

	st := spectralTypes[classIdx]

	star := models.Star{
		ID:           uuid.New().String(),
		Name:         fmt.Sprintf("Star-%d", index),
		SpectralType: spectralType,
		Mass:         randFloat(r, st.massRange[0], st.massRange[1]),
		Radius:       randFloat(r, st.radiusRange[0], st.radiusRange[1]),
		Temperature:  randInt(r, st.tempRange[0], st.tempRange[1]),
	}

	return star
}

// generatePlanet creates a realistic planet orbiting a star
func generatePlanet(r *rand.Rand, star models.Star, index int) models.Planet {
	// Orbital parameters
	semiMajorAxis := randFloat(r, 0.05, 50.0) // AU
	orbitalPeriod := 365.25 * math.Sqrt(semiMajorAxis*semiMajorAxis*semiMajorAxis/star.Mass)

	// Planet type determines mass and radius
	planetType := r.Float64()
	var mass, radius float64

	if planetType < 0.3 { // Rocky planet
		mass = randFloat(r, 0.1, 5.0)
		radius = randFloat(r, 0.3, 2.0)
	} else if planetType < 0.6 { // Ice giant
		mass = randFloat(r, 5.0, 20.0)
		radius = randFloat(r, 2.0, 4.0)
	} else { // Gas giant
		mass = randFloat(r, 20.0, 1000.0)
		radius = randFloat(r, 4.0, 15.0)
	}

	// Temperature decreases with distance
	surfaceTemp := int32(float64(star.Temperature) * math.Sqrt(star.Radius/(2.0*semiMajorAxis)))

	planet := models.Planet{
		ID:            uuid.New().String(),
		Name:          fmt.Sprintf("%s-Planet-%d", star.Name, index),
		OrbitalPeriod: orbitalPeriod,
		SemiMajorAxis: semiMajorAxis,
		Eccentricity:  randFloat(r, 0.0, 0.3),
		Mass:          mass,
		Radius:        radius,
		Atmosphere:    atmosphereTypes[r.Intn(len(atmosphereTypes))],
		SurfaceTemp:   surfaceTemp,
		HasRings:      r.Float64() < 0.2,
		HasMoons:      r.Float64() < 0.6,
		DiscoveryYear: randInt(r, 1990, 2024),
		StarID:        star.ID,
	}

	return planet
}

// generateExoplanet creates a realistic exoplanet
func generateExoplanet(r *rand.Rand, star models.Star, index int) models.Exoplanet {
	// Orbital parameters
	semiMajorAxis := randFloat(r, 0.01, 5.0) // AU (closer range for detectability)
	orbitalPeriod := 365.25 * math.Sqrt(semiMajorAxis*semiMajorAxis*semiMajorAxis/star.Mass)

	// Mass and radius
	mass := randFloat(r, 0.5, 500.0)
	radius := randFloat(r, 0.5, 12.0)

	// Temperature
	surfaceTemp := int32(float64(star.Temperature) * math.Sqrt(star.Radius/(2.0*semiMajorAxis)))

	// Host distance
	hostDistance := randFloat(r, 10.0, 10000.0)

	exoplanet := models.Exoplanet{
		ID:              uuid.New().String(),
		Name:            fmt.Sprintf("%s-Exo-%d", star.Name, index),
		OrbitalPeriod:   orbitalPeriod,
		SemiMajorAxis:   semiMajorAxis,
		Eccentricity:    randFloat(r, 0.0, 0.5),
		Mass:            mass,
		Radius:          radius,
		DetectionMethod: detectionMethods[r.Intn(len(detectionMethods))],
		HostDistance:    hostDistance,
		SurfaceTemp:     surfaceTemp,
		DiscoveryYear:   randInt(r, 1990, 2024),
		StarID:          star.ID,
	}

	return exoplanet
}

// GenerateAll generates all stellar data based on configuration
func GenerateAll(cfg Config) *GeneratedData {
	r := rand.New(rand.NewSource(cfg.Seed))

	data := &GeneratedData{
		Stars:      make([]models.Star, 0, cfg.NumStars),
		Planets:    make([]models.Planet, 0),
		Exoplanets: make([]models.Exoplanet, 0),
	}

	// Generate stars
	for i := 0; i < cfg.NumStars; i++ {
		star := generateStar(r, i+1)
		data.Stars = append(data.Stars, star)

		// Generate planets for this star
		numPlanets := r.Intn(cfg.PlanetsPerStar + 1) // 0 to PlanetsPerStar
		for j := 0; j < numPlanets; j++ {
			planet := generatePlanet(r, star, j+1)
			data.Planets = append(data.Planets, planet)
		}

		// Generate exoplanets for this star
		numExoplanets := r.Intn(cfg.ExoPerStar + 1) // 0 to ExoPerStar
		for j := 0; j < numExoplanets; j++ {
			exoplanet := generateExoplanet(r, star, j+1)
			data.Exoplanets = append(data.Exoplanets, exoplanet)
		}
	}

	return data
}
