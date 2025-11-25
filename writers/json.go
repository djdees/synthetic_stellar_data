package writers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"djdees/synthetic_stellar_data/generator"
)

// WriteJSON writes generated data to JSON files
func WriteJSON(data *generator.GeneratedData, outputDir string) error {
	// Write stars
	if err := writeJSONFile(data.Stars, filepath.Join(outputDir, "stars.json")); err != nil {
		return fmt.Errorf("failed to write stars JSON: %w", err)
	}

	// Write planets
	if err := writeJSONFile(data.Planets, filepath.Join(outputDir, "planets.json")); err != nil {
		return fmt.Errorf("failed to write planets JSON: %w", err)
	}

	// Write exoplanets
	if err := writeJSONFile(data.Exoplanets, filepath.Join(outputDir, "exoplanets.json")); err != nil {
		return fmt.Errorf("failed to write exoplanets JSON: %w", err)
	}

	return nil
}

func writeJSONFile(data interface{}, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}
