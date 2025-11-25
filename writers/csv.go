package writers

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"djdees/synthetic_stellar_data/generator"
	"djdees/synthetic_stellar_data/models"
)

// WriteCSV writes generated data to CSV files
func WriteCSV(data *generator.GeneratedData, outputDir string) error {
	// Write stars
	if err := writeStarsCSV(data.Stars, filepath.Join(outputDir, "stars.csv")); err != nil {
		return fmt.Errorf("failed to write stars CSV: %w", err)
	}

	// Write planets
	if err := writePlanetsCSV(data.Planets, filepath.Join(outputDir, "planets.csv")); err != nil {
		return fmt.Errorf("failed to write planets CSV: %w", err)
	}

	// Write exoplanets
	if err := writeExoplanetsCSV(data.Exoplanets, filepath.Join(outputDir, "exoplanets.csv")); err != nil {
		return fmt.Errorf("failed to write exoplanets CSV: %w", err)
	}

	return nil
}

func writeStarsCSV(stars []models.Star, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"ID", "Name", "SpectralType", "Mass", "Radius", "Temperature"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data
	for _, star := range stars {
		record := []string{
			star.ID,
			star.Name,
			star.SpectralType,
			fmt.Sprintf("%.6f", star.Mass),
			fmt.Sprintf("%.6f", star.Radius),
			fmt.Sprintf("%d", star.Temperature),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func writePlanetsCSV(planets []models.Planet, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"ID", "Name", "OrbitalPeriod", "SemiMajorAxis", "Eccentricity",
		"Mass", "Radius", "Atmosphere", "SurfaceTemp", "HasRings",
		"HasMoons", "DiscoveryYear", "StarID",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data
	for _, planet := range planets {
		record := []string{
			planet.ID,
			planet.Name,
			fmt.Sprintf("%.6f", planet.OrbitalPeriod),
			fmt.Sprintf("%.6f", planet.SemiMajorAxis),
			fmt.Sprintf("%.6f", planet.Eccentricity),
			fmt.Sprintf("%.6f", planet.Mass),
			fmt.Sprintf("%.6f", planet.Radius),
			planet.Atmosphere,
			fmt.Sprintf("%d", planet.SurfaceTemp),
			fmt.Sprintf("%t", planet.HasRings),
			fmt.Sprintf("%t", planet.HasMoons),
			fmt.Sprintf("%d", planet.DiscoveryYear),
			planet.StarID,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func writeExoplanetsCSV(exoplanets []models.Exoplanet, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"ID", "Name", "OrbitalPeriod", "SemiMajorAxis", "Eccentricity",
		"Mass", "Radius", "DetectionMethod", "HostDistance", "SurfaceTemp",
		"DiscoveryYear", "StarID",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data
	for _, exo := range exoplanets {
		record := []string{
			exo.ID,
			exo.Name,
			fmt.Sprintf("%.6f", exo.OrbitalPeriod),
			fmt.Sprintf("%.6f", exo.SemiMajorAxis),
			fmt.Sprintf("%.6f", exo.Eccentricity),
			fmt.Sprintf("%.6f", exo.Mass),
			fmt.Sprintf("%.6f", exo.Radius),
			exo.DetectionMethod,
			fmt.Sprintf("%.6f", exo.HostDistance),
			fmt.Sprintf("%d", exo.SurfaceTemp),
			fmt.Sprintf("%d", exo.DiscoveryYear),
			exo.StarID,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}
