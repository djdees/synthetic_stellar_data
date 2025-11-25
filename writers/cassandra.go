package writers

import (
	"fmt"
	"log"

	"djdees/synthetic_stellar_data/config"
	"djdees/synthetic_stellar_data/generator"
	"djdees/synthetic_stellar_data/models"
	gocql "github.com/apache/cassandra-gocql-driver/v2"
)

// WriteToCassandra writes generated data to Cassandra database
func WriteToCassandra(data *generator.GeneratedData, configFile string) error {
	// Load Cassandra configuration
	cfg, err := config.LoadCassandraConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create cluster configuration
	cluster := gocql.NewCluster(cfg.Hosts...)
	cluster.Keyspace = "system"

	if cfg.Username != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: cfg.Username,
			Password: cfg.Password,
		}
	}

	// Create session
	session, err := cluster.CreateSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Create keyspace
	if err := createKeyspace(session, cfg); err != nil {
		return err
	}

	// Switch to the keyspace
	session.Close()
	cluster.Keyspace = cfg.Keyspace
	session, err = cluster.CreateSession()
	if err != nil {
		return fmt.Errorf("failed to create session with keyspace: %w", err)
	}
	defer session.Close()

	// Create tables
	if err := createTables(session); err != nil {
		return err
	}

	// Insert data
	log.Println("Inserting stars...")
	if err := insertStars(session, data.Stars); err != nil {
		return err
	}

	log.Println("Inserting planets...")
	if err := insertPlanets(session, data.Planets); err != nil {
		return err
	}

	log.Println("Inserting exoplanets...")
	if err := insertExoplanets(session, data.Exoplanets); err != nil {
		return err
	}

	return nil
}

func createKeyspace(session *gocql.Session, cfg *config.CassandraConfig) error {
	query := fmt.Sprintf(`
		CREATE KEYSPACE IF NOT EXISTS %s
		WITH replication = {
			'class': '%s',
			'replication_factor': %d
		}
	`, cfg.Keyspace, cfg.ReplicationClass, cfg.ReplicationFactor)

	if err := session.Query(query).Exec(); err != nil {
		return fmt.Errorf("failed to create keyspace: %w", err)
	}

	log.Printf("Keyspace %s created or already exists\n", cfg.Keyspace)
	return nil
}

func createTables(session *gocql.Session) error {
	// Create stars table
	starsTable := `
		CREATE TABLE IF NOT EXISTS stars (
			id text PRIMARY KEY,
			name text,
			spectral_type text,
			mass double,
			radius double,
			temperature int
		)
	`
	if err := session.Query(starsTable).Exec(); err != nil {
		return fmt.Errorf("failed to create stars table: %w", err)
	}

	// Create planets table
	planetsTable := `
		CREATE TABLE IF NOT EXISTS planets (
			id text PRIMARY KEY,
			name text,
			orbital_period double,
			semi_major_axis double,
			eccentricity double,
			mass double,
			radius double,
			atmosphere text,
			surface_temp int,
			has_rings boolean,
			has_moons boolean,
			discovery_year int,
			star_id text
		)
	`
	if err := session.Query(planetsTable).Exec(); err != nil {
		return fmt.Errorf("failed to create planets table: %w", err)
	}

	// Create exoplanets table
	exoplanetsTable := `
		CREATE TABLE IF NOT EXISTS exoplanets (
			id text PRIMARY KEY,
			name text,
			orbital_period double,
			semi_major_axis double,
			eccentricity double,
			mass double,
			radius double,
			detection_method text,
			host_distance double,
			surface_temp int,
			discovery_year int,
			star_id text
		)
	`
	if err := session.Query(exoplanetsTable).Exec(); err != nil {
		return fmt.Errorf("failed to create exoplanets table: %w", err)
	}

	log.Println("Tables created or already exist")
	return nil
}

func insertStars(session *gocql.Session, stars []models.Star) error {
	query := `
		INSERT INTO stars (id, name, spectral_type, mass, radius, temperature)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	// Use individual inserts instead of batching for Apache driver v2
	for _, star := range stars {
		if err := session.Query(query,
			star.ID,
			star.Name,
			star.SpectralType,
			star.Mass,
			star.Radius,
			star.Temperature,
		).Exec(); err != nil {
			return fmt.Errorf("failed to insert star %s: %w", star.Name, err)
		}
	}

	return nil
}

func insertPlanets(session *gocql.Session, planets []models.Planet) error {
	query := `
		INSERT INTO planets (id, name, orbital_period, semi_major_axis, eccentricity,
			mass, radius, atmosphere, surface_temp, has_rings, has_moons, discovery_year, star_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	// Use individual inserts
	for _, planet := range planets {
		if err := session.Query(query,
			planet.ID,
			planet.Name,
			planet.OrbitalPeriod,
			planet.SemiMajorAxis,
			planet.Eccentricity,
			planet.Mass,
			planet.Radius,
			planet.Atmosphere,
			planet.SurfaceTemp,
			planet.HasRings,
			planet.HasMoons,
			planet.DiscoveryYear,
			planet.StarID,
		).Exec(); err != nil {
			return fmt.Errorf("failed to insert planet %s: %w", planet.Name, err)
		}
	}

	return nil
}

func insertExoplanets(session *gocql.Session, exoplanets []models.Exoplanet) error {
	query := `
		INSERT INTO exoplanets (id, name, orbital_period, semi_major_axis, eccentricity,
			mass, radius, detection_method, host_distance, surface_temp, discovery_year, star_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	// Use individual inserts
	for _, exo := range exoplanets {
		if err := session.Query(query,
			exo.ID,
			exo.Name,
			exo.OrbitalPeriod,
			exo.SemiMajorAxis,
			exo.Eccentricity,
			exo.Mass,
			exo.Radius,
			exo.DetectionMethod,
			exo.HostDistance,
			exo.SurfaceTemp,
			exo.DiscoveryYear,
			exo.StarID,
		).Exec(); err != nil {
			return fmt.Errorf("failed to insert exoplanet %s: %w", exo.Name, err)
		}
	}

	return nil
}
