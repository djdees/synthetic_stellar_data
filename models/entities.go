package models

// Star represents a stellar object with physical characteristics
type Star struct {
	ID           string  // UUID string
	Name         string  // Star name
	SpectralType string  // Spectral classification (e.g., "G2V", "M5V")
	Mass         float64 // Mass in solar masses
	Radius       float64 // Radius in solar radii
	Temperature  int32   // Surface temperature in Kelvin
}

// Planet represents a planet orbiting a star
type Planet struct {
	ID            string  // UUID string
	Name          string  // Planet name
	OrbitalPeriod float64 // Orbital period in Earth days
	SemiMajorAxis float64 // Semi-major axis in AU
	Eccentricity  float64 // Orbital eccentricity (0-1)
	Mass          float64 // Mass in Earth masses
	Radius        float64 // Radius in Earth radii
	Atmosphere    string  // Atmospheric composition description
	SurfaceTemp   int32   // Surface temperature in Kelvin
	HasRings      bool    // Whether the planet has rings
	HasMoons      bool    // Whether the planet has moons
	DiscoveryYear int32   // Year of discovery
	StarID        string  // Foreign key to parent Star
}

// Exoplanet represents an exoplanet orbiting a distant star
type Exoplanet struct {
	ID              string  // UUID string
	Name            string  // Exoplanet name
	OrbitalPeriod   float64 // Orbital period in Earth days
	SemiMajorAxis   float64 // Semi-major axis in AU
	Eccentricity    float64 // Orbital eccentricity (0-1)
	Mass            float64 // Mass in Earth masses
	Radius          float64 // Radius in Earth radii
	DetectionMethod string  // Method used to detect the exoplanet
	HostDistance    float64 // Distance to host star in light years
	SurfaceTemp     int32   // Surface temperature in Kelvin
	DiscoveryYear   int32   // Year of discovery
	StarID          string  // Foreign key to parent Star
}
