# Stellargen - Synthetic Stellar Data Generator

A comprehensive GoLang application that generates synthetic stellar data (stars, planets, and exoplanets) for database load testing, data analysis, and benchmarking purposes.

## Features

- **Realistic Data Generation**: Creates scientifically plausible stellar systems with proper physical relationships
- **Multiple Output Formats**: Supports CSV, JSON, Parquet, and direct Cassandra insertion
- **Reproducible**: Seeded random generation ensures consistent results
- **Scalable**: Generate from hundreds to millions of records
- **Modular Architecture**: Clean separation of concerns with distinct packages
- **Comprehensive Testing**: Unit tests and benchmarks included

## Installation

### Prerequisites

- Go 1.21 or higher
- (Optional) Apache Cassandra for database output

### Build from Source

```bash
git clone <repository-url>
cd stellargen
make install-deps
make build
```

The binary will be created in the `bin/` directory.

## Usage

### Basic Command Structure

```bash
./stellargen [flags]
```

### Command-Line Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--num-stars` | int | 100 | Number of stars to generate |
| `--planets-per-star` | int | 8 | Maximum planets per star (0-15) |
| `--exo-per-star` | int | 5 | Maximum exoplanets per star (0-8) |
| `--output-format` | string | csv | Output format: csv, json, parquet, cassandra |
| `--output-dir` | string | output | Output directory for files |
| `--config` | string | | YAML config file for Cassandra |
| `--seed` | int64 | 0 | Random seed (0 for time-based) |
| `--dry-run` | bool | false | Generate data without writing output |

### Examples

#### Generate CSV Files

```bash
# Default settings (100 stars)
./stellargen

# Custom number of stars
./stellargen --num-stars=1000 --output-dir=./data

# Reproducible generation with seed
./stellargen --num-stars=500 --seed=12345
```

#### Generate JSON Files

```bash
./stellargen --num-stars=500 --output-format=json
```

#### Generate Parquet Files

```bash
./stellargen --num-stars=10000 --output-format=parquet
```

#### Insert into Cassandra

```bash
./stellargen --num-stars=1000 --output-format=cassandra --config=examples/config.yaml
```

#### Dry Run (Test Generation)

```bash
./stellargen --num-stars=10000 --dry-run
```

## Data Models

### Star

| Field | Type | Description |
|-------|------|-------------|
| ID | UUID | Unique identifier |
| Name | string | Star name |
| SpectralType | string | Spectral classification (e.g., "G2V", "M5V") |
| Mass | float64 | Mass in solar masses |
| Radius | float64 | Radius in solar radii |
| Temperature | int32 | Surface temperature in Kelvin |

### Planet

| Field | Type | Description |
|-------|------|-------------|
| ID | UUID | Unique identifier |
| Name | string | Planet name |
| OrbitalPeriod | float64 | Orbital period in Earth days |
| SemiMajorAxis | float64 | Semi-major axis in AU |
| Eccentricity | float64 | Orbital eccentricity (0-1) |
| Mass | float64 | Mass in Earth masses |
| Radius | float64 | Radius in Earth radii |
| Atmosphere | string | Atmospheric composition |
| SurfaceTemp | int32 | Surface temperature in Kelvin |
| HasRings | bool | Whether the planet has rings |
| HasMoons | bool | Whether the planet has moons |
| DiscoveryYear | int32 | Year of discovery (1990-2024) |
| StarID | string | Parent star UUID |

### Exoplanet

| Field | Type | Description |
|-------|------|-------------|
| ID | UUID | Unique identifier |
| Name | string | Exoplanet name |
| OrbitalPeriod | float64 | Orbital period in Earth days |
| SemiMajorAxis | float64 | Semi-major axis in AU |
| Eccentricity | float64 | Orbital eccentricity (0-1) |
| Mass | float64 | Mass in Earth masses |
| Radius | float64 | Radius in Earth radii |
| DetectionMethod | string | Detection method used |
| HostDistance | float64 | Distance to host star in light years |
| SurfaceTemp | int32 | Surface temperature in Kelvin |
| DiscoveryYear | int32 | Year of discovery (1990-2024) |
| StarID | string | Parent star UUID |

## Output Formats

### CSV

Three files are created:
- `stars.csv` - Star data with headers
- `planets.csv` - Planet data with headers
- `exoplanets.csv` - Exoplanet data with headers

### JSON

Three JSON files with pretty-printed output:
- `stars.json`
- `planets.json`
- `exoplanets.json`

### Parquet

Three Parquet files with columnar storage:
- `stars.parquet`
- `planets.parquet`
- `exoplanets.parquet`

### Cassandra

Data is inserted directly into Cassandra tables:
- `stars` table
- `planets` table
- `exoplanets` table

## Cassandra Configuration

Create a YAML configuration file (see `examples/config.yaml`):

```yaml
hosts:
  - 127.0.0.1
keyspace: stellargen
consistency: QUORUM
replication_class: SimpleStrategy
replication_factor: 1
username: ""
password: ""
```

## Development

### Project Structure

```
stellargen/
├── main.go                 # Application entry point
├── go.mod                  # Go module definition
├── Makefile               # Build automation
├── README.md              # This file
├── config/                # Configuration handling
│   ├── config.go          # Cassandra config
│   └── flags.go           # Command-line flags
├── models/                # Data models
│   └── entities.go        # Star, Planet, Exoplanet
├── generator/             # Data generation logic
│   ├── config.go          # Generator config
│   └── generate.go        # Generation functions
├── output/                # Output writers
│   ├── csv.go             # CSV output
│   ├── json.go            # JSON output
│   ├── parquet.go         # Parquet output
│   └── cassandra.go       # Cassandra output
├── tests/                 # Test suite
│   └── generator_test.go  # Unit tests
└── examples/              # Example configurations
    ├── config.yaml        # Cassandra config example
    └── stress-test/       # Stress test profiles
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run benchmarks
go test -bench=. ./tests/
```

### Make Commands

```bash
make help              # Show all available commands
make build             # Compile the application
make test              # Run tests
make test-coverage     # Run tests with coverage
make clean             # Remove build artifacts
make run               # Build and run with defaults
make run-csv           # Generate CSV files
make run-json          # Generate JSON files
make run-parquet       # Generate Parquet files
make run-cassandra     # Insert into Cassandra
make install-deps      # Install dependencies
make fmt               # Format code
make vet               # Run go vet
```

## Data Generation Details

### Spectral Types

Stars are generated with realistic spectral classifications:
- **O stars**: 30,000-50,000 K, 16-90 solar masses (very rare)
- **B stars**: 10,000-30,000 K, 2.1-16 solar masses
- **A stars**: 7,500-10,000 K, 1.4-2.1 solar masses
- **F stars**: 6,000-7,500 K, 1.04-1.4 solar masses
- **G stars**: 5,200-6,000 K, 0.8-1.04 solar masses (like our Sun)
- **K stars**: 3,700-5,200 K, 0.45-0.8 solar masses
- **M stars**: 2,400-3,700 K, 0.08-0.45 solar masses (most common)

### Planetary Types

Three main categories:
- **Rocky planets**: 0.1-5 Earth masses, 0.3-2 Earth radii
- **Ice giants**: 5-20 Earth masses, 2-4 Earth radii
- **Gas giants**: 20-1000 Earth masses, 4-15 Earth radii

### Detection Methods

Exoplanets use realistic detection methods:
- Transit
- Radial Velocity
- Direct Imaging
- Gravitational Microlensing
- Astrometry
- Transit Timing Variation

## Performance

Typical generation speeds (on modern hardware):
- 1,000 stars: < 100ms
- 10,000 stars: < 1 second
- 100,000 stars: < 10 seconds

## License

[Specify your license here]

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

For issues, questions, or contributions, please open an issue on the repository.
