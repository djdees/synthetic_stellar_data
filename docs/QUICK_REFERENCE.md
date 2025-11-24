# Stellargen Quick Reference

## Installation & Setup

```bash
git clone <repo-url> && cd stellargen
make install-deps
make build
```

## Common Commands

```bash
# Basic generation (CSV, 100 stars)
./bin/stellargen

# Specify number of stars
./bin/stellargen --num-stars=1000

# Different output formats
./bin/stellargen --output-format=json
./bin/stellargen --output-format=parquet
./bin/stellargen --output-format=cassandra --config=examples/config.yaml

# Reproducible with seed
./bin/stellargen --seed=12345

# Dry run (no output)
./bin/stellargen --num-stars=10000 --dry-run

# Custom output directory
./bin/stellargen --output-dir=./mydata

# Control planetary distribution
./bin/stellargen --planets-per-star=10 --exo-per-star=3
```

## Make Targets

```bash
make build              # Compile binary
make test               # Run tests
make test-coverage      # Run tests with coverage
make clean              # Remove artifacts
make run                # Build and run
make run-csv            # Generate CSV (1000 stars)
make run-json           # Generate JSON (500 stars)
make run-parquet        # Generate Parquet (10K stars)
make run-cassandra      # Insert to Cassandra
make install-deps       # Install Go dependencies
make fmt                # Format code
make vet                # Run go vet
make all                # Clean, format, test, build
```

## Command-Line Flags

| Flag | Default | Range/Options | Description |
|------|---------|---------------|-------------|
| `--num-stars` | 100 | 1+ | Number of stars |
| `--planets-per-star` | 8 | 0-15 | Max planets per star |
| `--exo-per-star` | 5 | 0-8 | Max exoplanets per star |
| `--output-format` | csv | csv, json, parquet, cassandra | Output format |
| `--output-dir` | output | any path | Output directory |
| `--config` | "" | path to yaml | Cassandra config |
| `--seed` | 0 | int64 | Random seed (0=time) |
| `--dry-run` | false | bool | Generate without writing |

## Data Models Summary

### Star
- ID (UUID), Name, SpectralType (e.g., "G2V")
- Mass (solar masses), Radius (solar radii), Temperature (K)

### Planet
- ID (UUID), Name, OrbitalPeriod (days), SemiMajorAxis (AU)
- Eccentricity (0-1), Mass (Earth masses), Radius (Earth radii)
- Atmosphere, SurfaceTemp (K), HasRings, HasMoons
- DiscoveryYear (1990-2024), StarID (FK)

### Exoplanet
- ID (UUID), Name, OrbitalPeriod (days), SemiMajorAxis (AU)
- Eccentricity (0-1), Mass (Earth masses), Radius (Earth radii)
- DetectionMethod, HostDistance (ly), SurfaceTemp (K)
- DiscoveryYear (1990-2024), StarID (FK)

## Output Files

### CSV
- `output/stars.csv`
- `output/planets.csv`
- `output/exoplanets.csv`

### JSON
- `output/stars.json`
- `output/planets.json`
- `output/exoplanets.json`

### Parquet
- `output/stars.parquet`
- `output/planets.parquet`
- `output/exoplanets.parquet`

### Cassandra
- Tables: `stars`, `planets`, `exoplanets`
- Keyspace: from config.yaml

## Cassandra Quick Setup

```bash
# Start Cassandra (Docker)
docker run --name cassandra -p 9042:9042 -d cassandra:latest

# Wait for startup (30-60 seconds)
docker logs -f cassandra

# Create config (examples/config.yaml already exists)
# Run stellargen
./bin/stellargen --num-stars=1000 --output-format=cassandra --config=examples/config.yaml

# Query data
docker exec -it cassandra cqlsh
cqlsh> USE stellargen;
cqlsh:stellargen> SELECT COUNT(*) FROM stars;
```

## Testing

```bash
# Unit tests
go test ./tests/... -v

# With coverage
go test ./tests/... -cover

# Benchmarks
go test ./tests/... -bench=.

# Specific test
go test ./tests/... -run TestGenerateAllBasic
```

## Performance Estimates

| Stars | Generation Time | CSV Write | Parquet Write | Memory |
|-------|----------------|-----------|---------------|--------|
| 1K | <100ms | ~50ms | ~100ms | ~10MB |
| 10K | <1s | ~500ms | ~1s | ~100MB |
| 100K | ~10s | ~5s | ~10s | ~1GB |
| 1M | ~2min | ~50s | ~2min | ~10GB |

## Spectral Types

| Class | Temp (K) | Mass (M☉) | Color | Frequency |
|-------|----------|-----------|-------|-----------|
| O | 30K-50K | 16-90 | Blue | Very Rare |
| B | 10K-30K | 2.1-16 | Blue-white | Rare |
| A | 7.5K-10K | 1.4-2.1 | White | Uncommon |
| F | 6K-7.5K | 1.04-1.4 | Yellow-white | Common |
| G | 5.2K-6K | 0.8-1.04 | Yellow | Common |
| K | 3.7K-5.2K | 0.45-0.8 | Orange | Very Common |
| M | 2.4K-3.7K | 0.08-0.45 | Red | Most Common |

## Detection Methods

- Transit
- Radial Velocity
- Direct Imaging
- Gravitational Microlensing
- Astrometry
- Transit Timing Variation

## Atmospheric Types

- H2/He dominant (gas giants)
- N2/O2 dominant (Earth-like)
- CO2 dominant (Venus-like)
- Thin atmosphere
- No atmosphere
- H2/He with methane
- Sulfuric compounds
- Water vapor rich

## Common Issues

| Problem | Solution |
|---------|----------|
| Can't build | Run `make install-deps` |
| Permission denied | `chmod +x bin/stellargen` or `mkdir -p output` |
| Cassandra connection fails | Check Docker/service running: `docker ps` |
| Out of memory | Reduce `--num-stars` or use `--output-format=cassandra` |
| Different results | Use same `--seed` value |

## File Structure

```
stellargen/
├── main.go              # Entry point
├── go.mod               # Dependencies
├── Makefile             # Build automation
├── config/              # Config parsing
├── models/              # Data models
├── generator/           # Data generation
├── output/              # Output writers
├── tests/               # Test suite
└── examples/            # Config examples
```

## Integration Examples

### Python (Pandas)
```python
import pandas as pd
df = pd.read_parquet('output/stars.parquet')
print(df.describe())
```

### SQL Import
```sql
-- PostgreSQL
COPY stars FROM '/path/to/stars.csv' CSV HEADER;
```

### MongoDB
```bash
mongoimport --db astronomy --collection stars --file output/stars.json --jsonArray
```

## Environment Variables

None required. All configuration via flags or config file.

## Exit Codes

- 0: Success
- 1: Error (see error message)

## Support

- GitHub Issues: [repo-url]/issues
- Documentation: README.md, USAGE_GUIDE.md
- Examples: examples/ directory
