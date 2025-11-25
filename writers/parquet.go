package writers

import (
	"fmt"
	"path/filepath"

	"djdees/synthetic_stellar_data/generator"
	"djdees/synthetic_stellar_data/models"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/writer"
)

// Parquet-tagged structs for efficient columnar storage
type StarParquet struct {
	ID           string  `parquet:"name=id, type=BYTE_ARRAY, convertedtype=UTF8"`
	Name         string  `parquet:"name=name, type=BYTE_ARRAY, convertedtype=UTF8"`
	SpectralType string  `parquet:"name=spectral_type, type=BYTE_ARRAY, convertedtype=UTF8"`
	Mass         float64 `parquet:"name=mass, type=DOUBLE"`
	Radius       float64 `parquet:"name=radius, type=DOUBLE"`
	Temperature  int32   `parquet:"name=temperature, type=INT32"`
}

type PlanetParquet struct {
	ID            string  `parquet:"name=id, type=BYTE_ARRAY, convertedtype=UTF8"`
	Name          string  `parquet:"name=name, type=BYTE_ARRAY, convertedtype=UTF8"`
	OrbitalPeriod float64 `parquet:"name=orbital_period, type=DOUBLE"`
	SemiMajorAxis float64 `parquet:"name=semi_major_axis, type=DOUBLE"`
	Eccentricity  float64 `parquet:"name=eccentricity, type=DOUBLE"`
	Mass          float64 `parquet:"name=mass, type=DOUBLE"`
	Radius        float64 `parquet:"name=radius, type=DOUBLE"`
	Atmosphere    string  `parquet:"name=atmosphere, type=BYTE_ARRAY, convertedtype=UTF8"`
	SurfaceTemp   int32   `parquet:"name=surface_temp, type=INT32"`
	HasRings      bool    `parquet:"name=has_rings, type=BOOLEAN"`
	HasMoons      bool    `parquet:"name=has_moons, type=BOOLEAN"`
	DiscoveryYear int32   `parquet:"name=discovery_year, type=INT32"`
	StarID        string  `parquet:"name=star_id, type=BYTE_ARRAY, convertedtype=UTF8"`
}

type ExoplanetParquet struct {
	ID              string  `parquet:"name=id, type=BYTE_ARRAY, convertedtype=UTF8"`
	Name            string  `parquet:"name=name, type=BYTE_ARRAY, convertedtype=UTF8"`
	OrbitalPeriod   float64 `parquet:"name=orbital_period, type=DOUBLE"`
	SemiMajorAxis   float64 `parquet:"name=semi_major_axis, type=DOUBLE"`
	Eccentricity    float64 `parquet:"name=eccentricity, type=DOUBLE"`
	Mass            float64 `parquet:"name=mass, type=DOUBLE"`
	Radius          float64 `parquet:"name=radius, type=DOUBLE"`
	DetectionMethod string  `parquet:"name=detection_method, type=BYTE_ARRAY, convertedtype=UTF8"`
	HostDistance    float64 `parquet:"name=host_distance, type=DOUBLE"`
	SurfaceTemp     int32   `parquet:"name=surface_temp, type=INT32"`
	DiscoveryYear   int32   `parquet:"name=discovery_year, type=INT32"`
	StarID          string  `parquet:"name=star_id, type=BYTE_ARRAY, convertedtype=UTF8"`
}

// WriteParquet writes generated data to Parquet files
func WriteParquet(data *generator.GeneratedData, outputDir string) error {
	// Write stars
	if err := writeStarsParquet(data.Stars, filepath.Join(outputDir, "stars.parquet")); err != nil {
		return fmt.Errorf("failed to write stars parquet: %w", err)
	}

	// Write planets
	if err := writePlanetsParquet(data.Planets, filepath.Join(outputDir, "planets.parquet")); err != nil {
		return fmt.Errorf("failed to write planets parquet: %w", err)
	}

	// Write exoplanets
	if err := writeExoplanetsParquet(data.Exoplanets, filepath.Join(outputDir, "exoplanets.parquet")); err != nil {
		return fmt.Errorf("failed to write exoplanets parquet: %w", err)
	}

	return nil
}

func writeStarsParquet(stars []models.Star, filename string) error {
	fw, err := local.NewLocalFileWriter(filename)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := fw.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	pw, err := writer.NewParquetWriter(fw, new(StarParquet), 4)
	if err != nil {
		return err
	}

	for _, star := range stars {
		s := StarParquet{
			ID:           star.ID,
			Name:         star.Name,
			SpectralType: star.SpectralType,
			Mass:         star.Mass,
			Radius:       star.Radius,
			Temperature:  star.Temperature,
		}
		if err := pw.Write(s); err != nil {
			return err
		}
	}

	if err := pw.WriteStop(); err != nil {
		return err
	}

	return nil
}

func writePlanetsParquet(planets []models.Planet, filename string) error {
	fw, err := local.NewLocalFileWriter(filename)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := fw.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	pw, err := writer.NewParquetWriter(fw, new(PlanetParquet), 4)
	if err != nil {
		return err
	}

	for _, planet := range planets {
		p := PlanetParquet{
			ID:            planet.ID,
			Name:          planet.Name,
			OrbitalPeriod: planet.OrbitalPeriod,
			SemiMajorAxis: planet.SemiMajorAxis,
			Eccentricity:  planet.Eccentricity,
			Mass:          planet.Mass,
			Radius:        planet.Radius,
			Atmosphere:    planet.Atmosphere,
			SurfaceTemp:   planet.SurfaceTemp,
			HasRings:      planet.HasRings,
			HasMoons:      planet.HasMoons,
			DiscoveryYear: planet.DiscoveryYear,
			StarID:        planet.StarID,
		}
		if err := pw.Write(p); err != nil {
			return err
		}
	}

	if err := pw.WriteStop(); err != nil {
		return err
	}

	return nil
}

func writeExoplanetsParquet(exoplanets []models.Exoplanet, filename string) error {
	fw, err := local.NewLocalFileWriter(filename)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := fw.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	pw, err := writer.NewParquetWriter(fw, new(ExoplanetParquet), 4)
	if err != nil {
		return err
	}

	for _, exo := range exoplanets {
		e := ExoplanetParquet{
			ID:              exo.ID,
			Name:            exo.Name,
			OrbitalPeriod:   exo.OrbitalPeriod,
			SemiMajorAxis:   exo.SemiMajorAxis,
			Eccentricity:    exo.Eccentricity,
			Mass:            exo.Mass,
			Radius:          exo.Radius,
			DetectionMethod: exo.DetectionMethod,
			HostDistance:    exo.HostDistance,
			SurfaceTemp:     exo.SurfaceTemp,
			DiscoveryYear:   exo.DiscoveryYear,
			StarID:          exo.StarID,
		}
		if err := pw.Write(e); err != nil {
			return err
		}
	}

	if err := pw.WriteStop(); err != nil {
		return err
	}

	return nil
}
