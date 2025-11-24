# Stellargen Usage Guide

This guide provides detailed examples and use cases for the Stellargen application.

## Table of Contents

1. [Quick Start](#quick-start)
2. [Common Use Cases](#common-use-cases)
3. [Output Format Examples](#output-format-examples)
4. [Cassandra Setup](#cassandra-setup)
5. [Performance Tuning](#performance-tuning)
6. [Troubleshooting](#troubleshooting)

## Quick Start

### Build and Run

```bash
# Clone and build
git clone <repository-url>
cd stellargen
make build

# Generate sample data
./bin/stellargen --num-stars=100
```

Your data will be in the `output/` directory as CSV files.

### First Steps

```bash
# 1. Generate a small dataset to explore
./bin/stellargen --num-stars=10 --dry-run

# 2. Generate actual CSV files
./bin/stellargen --num-stars=100

# 3. View the results
head output/stars.csv
head output/planets.csv
head output/exoplanets.csv
```

## Common Use Cases

### 1. Database Load Testing

Generate large datasets for testing database performance:

```bash
# Small test (1K stars, ~5K planets, ~3K exoplanets)
./bin/stellargen --num-stars=1000 --output-format=csv

# Medium load (10K stars)
./bin/stellargen --num-stars=10000 --output-format=csv

# Large load (100K stars)
./bin/stellargen --num-stars=100000 --output-format=parquet

# Directly to Cassandra for load testing
./bin/stellargen --num-stars=50000 \
  --output-format=cassandra \
  --config=examples/config.yaml
```

### 2. Data Analysis and Visualization

Generate data for analytics tools:

```bash
# JSON for web applications
./bin/stellargen --num-stars=1000 \
  --output-format=json \
  --output-dir=./web/data

# Parquet for big data analytics (Spark, Pandas)
./bin/stellargen --num-stars=50000 \
  --output-format=parquet \
  --output-dir=./analytics/data
```

### 3. Reproducible Datasets

Create consistent datasets across environments:

```bash
# Development dataset (always same data)
./bin/stellargen --num-stars=500 \
  --seed=12345 \
  --output-dir=./dev-data

# Testing dataset (different seed)
./bin/stellargen --num-stars=500 \
  --seed=67890 \
  --output-dir=./test-data

# These will generate identical results every time:
./bin/stellargen --num-stars=100 --seed=99999
./bin/stellargen --num-stars=100 --seed=99999
```

### 4. Specific Planetary System Distributions

Adjust the mix of planets and exoplanets:

```bash
# Systems with many planets (like our solar system)
./bin/stellargen --num-stars=1000 \
  --planets-per-star=15 \
  --exo-per-star=0

# Exoplanet-focused dataset
./bin/stellargen --num-stars=5000 \
  --planets-per-star=0 \
  --exo-per-star=8

# Balanced systems
./bin/stellargen --num-stars=1000 \
  --planets-per-star=8 \
  --exo-per-star=5
```

### 5. Performance Benchmarking

Test generation and output performance:

```bash
# Benchmark generation only (no I/O)
time ./bin/stellargen --num-stars=10000 --dry-run

# Benchmark CSV writing
time ./bin/stellargen --num-stars=10000 --output-format=csv

# Benchmark Parquet writing
time ./bin/stellargen --num-stars=10000 --output-format=parquet

# Benchmark database insertion
time ./bin/stellargen --num-stars=10000 \
  --output-format=cassandra \
  --config=examples/config.yaml
```

## Output Format Examples

### CSV Output

```csv
# stars.csv
ID,Name,SpectralType,Mass,Radius,Temperature
550e8400-e29b-41d4-a716-446655440000,Star-1,G2V,0.985432,1.023456,5778
...

# planets.csv
ID,Name,OrbitalPeriod,SemiMajorAxis,Eccentricity,Mass,Radius,Atmosphere,SurfaceTemp,HasRings,HasMoons,DiscoveryYear,StarID
...

# exoplanets.csv
ID,Name,OrbitalPeriod,SemiMajorAxis,Eccentricity,Mass,Radius,DetectionMethod,HostDistance,SurfaceTemp,DiscoveryYear,StarID
...
```

**Use cases:**
- Spreadsheet analysis (Excel, Google Sheets)
- SQL database imports
- Simple data inspection

### JSON Output

```json
// stars.json
[
  {
    "ID": "550e8400-e29b-41d4-a716-446655440000",
    "Name": "Star-1",
    "SpectralType": "G2V",
    "Mass": 0.985432,
    "Radius": 1.023456,
    "Temperature": 5778
  },
  ...
]
```

**Use cases:**
- Web applications and APIs
- NoSQL databases (MongoDB, Elasticsearch)
- JavaScript/Python data processing

### Parquet Output

Binary columnar format, best viewed with specialized tools:

```bash
# View Parquet schema
parquet-tools schema output/stars.parquet

# View sample data
parquet-tools head output/stars.parquet

# Use in Python with pandas
python -c "import pandas as pd; df = pd.read_parquet('output/stars.parquet'); print(df.head())"
```

**Use cases:**
- Big data processing (Apache Spark, Hadoop)
- Data warehouses (Snowflake, BigQuery)
- Efficient storage and analytics

## Cassandra Setup

### 1. Install Cassandra

**Using Docker:**
```bash
docker run --name cassandra -p 9042:9042 -d cassandra:latest
```

**Using package manager:**
```bash
# Ubuntu/Debian
sudo apt install cassandra

# macOS
brew install cassandra
```

### 2. Verify Cassandra is Running

```bash
# Check status
docker exec -it cassandra nodetool status

# Or connect with cqlsh
docker exec -it cassandra cqlsh
```

### 3. Configure stellargen

Edit `examples/config.yaml`:

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

### 4. Load Data

```bash
./bin/stellargen --num-stars=10000 \
  --output-format=cassandra \
  --config=examples/config.yaml
```

### 5. Query Data

```cql
cqlsh> USE stellargen;
cqlsh:stellargen> SELECT COUNT(*) FROM stars;
cqlsh:stellargen> SELECT * FROM stars LIMIT 10;
cqlsh:stellargen> SELECT name, mass, temperature FROM stars WHERE spectral_type = 'G2V' ALLOW FILTERING;
```

## Performance Tuning

### Generation Performance

The generation speed depends on the number of entities:

| Stars | Estimated Time | Memory Usage |
|-------|----------------|--------------|
| 1K | < 100ms | ~10 MB |
| 10K | < 1s | ~100 MB |
| 100K | < 10s | ~1 GB |
| 1M | < 2 min | ~10 GB |

**Tips:**
- Use `--dry-run` to test generation speed without I/O
- Large datasets benefit from Parquet format (faster write)
- Cassandra batches are optimized for 100 records at a time

### Output Performance

**CSV:**
- Fast for small datasets (< 100K stars)
- Human-readable, easy to debug
- Larger file size

**JSON:**
- Similar to CSV performance
- Better for web applications
- Structured format

**Parquet:**
- Slower write, much faster read
- Columnar storage (better compression)
- Best for analytics workloads

**Cassandra:**
- Direct insertion, no intermediate files
- Batched for efficiency
- Network latency can affect speed

### Memory Optimization

For very large datasets:

```bash
# Generate in chunks
for i in {1..10}; do
  ./bin/stellargen --num-stars=100000 \
    --seed=$RANDOM \
    --output-dir=output/batch_$i
done
```

## Troubleshooting

### Issue: "Failed to create output directory"

**Solution:**
```bash
# Ensure you have write permissions
mkdir -p output
chmod 755 output
```

### Issue: "Cassandra connection failed"

**Solutions:**
```bash
# Check Cassandra is running
docker ps | grep cassandra

# Check configuration file exists
cat examples/config.yaml

# Test connection manually
cqlsh 127.0.0.1 9042
```

### Issue: "Out of memory"

**Solutions:**
```bash
# Reduce dataset size
./bin/stellargen --num-stars=10000  # instead of 1000000

# Use Cassandra output (streams data)
./bin/stellargen --num-stars=100000 --output-format=cassandra

# Increase system memory limits
ulimit -v unlimited
```

### Issue: "Generation is too slow"

**Solutions:**
```bash
# Use fixed seed (avoids random source contention)
./bin/stellargen --seed=12345

# Reduce complexity
./bin/stellargen --planets-per-star=5 --exo-per-star=3

# Use dry-run to isolate generation vs I/O
./bin/stellargen --num-stars=100000 --dry-run
```

### Issue: "Invalid spectral type generated"

This shouldn't happen, but if it does:
```bash
# Report the seed used
./bin/stellargen --seed=<problem-seed> --dry-run

# Use a different seed
./bin/stellargen --seed=12345
```

## Advanced Usage

### Custom Pipeline

```bash
# Generate, transform, and load
./bin/stellargen --num-stars=10000 --output-format=json

# Process with jq
cat output/stars.json | jq '.[] | select(.Temperature > 10000)' > hot_stars.json

# Import to database
mongoimport --db astronomy --collection stars --file hot_stars.json --jsonArray
```

### Parallel Generation

```bash
# Generate multiple datasets in parallel
parallel -j 4 "./bin/stellargen --num-stars=25000 --seed={} --output-dir=output/batch_{}" ::: {1..4}

# Merge CSV files
cat output/batch_*/stars.csv > all_stars.csv
```

### Integration with Data Pipelines

```python
# Python example with pandas
import pandas as pd
import subprocess

# Generate data
subprocess.run([
    './bin/stellargen',
    '--num-stars=10000',
    '--output-format=parquet'
])

# Load and analyze
stars = pd.read_parquet('output/stars.parquet')
planets = pd.read_parquet('output/planets.parquet')

print(f"Total stars: {len(stars)}")
print(f"Total planets: {len(planets)}")
print(f"Average planets per star: {len(planets) / len(stars):.2f}")
```

## Best Practices

1. **Always use seeds for reproducibility** in testing environments
2. **Start small** when testing new configurations
3. **Use Parquet** for large analytical datasets
4. **Use Cassandra output** for database load testing
5. **Monitor memory** with large datasets
6. **Version your data** by including seed in output directory name:
   ```bash
   ./bin/stellargen --seed=12345 --output-dir=output/seed_12345
   ```

