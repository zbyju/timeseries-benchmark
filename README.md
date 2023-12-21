# Time series database Benchmark

A benchmark suite to understand performance of selected time series databases (and also general purpose databases that can be used for this purpose) in real-world scenarios.

Benchmarks:

- Inserting data
- Querying last day's data
- Aggregating average values over all time

## Running benchmarks

```sh
# Launch all databases
docker-compose up

# Initialize InfluxDB on http://localhost:8086

# Run the benchmark
go test -bench=. -benchtime=1x
```

You can reset all databases by deleting the container along the volume.
