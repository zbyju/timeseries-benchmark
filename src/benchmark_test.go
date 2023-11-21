package main

import (
	"benchmark/simulator"
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"
	"time"

	// Import database drivers
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	pgx "github.com/jackc/pgx/v4"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

func printSlice(slice []float64) string {
	resultStrs := make([]string, len(slice))
	for i, val := range slice {
		resultStrs[i] = strconv.FormatFloat(val, 'f', -1, 64)
	}
	return strings.Join(resultStrs, "\n")
}

type Sample struct {
	StationId int
	Timestamp time.Time
	Voltage   float64
	Daily     float64
	Heating   float64
	Cooling   float64
}

func generateData() [][][]Sample {
	NumSamples := 144
	NumSimulations := 10
	NumStations := 100
	var result [][][]Sample

	for stationId := 0; stationId < NumStations; stationId++ {
		var simulations [][]Sample
		date := time.Now()
		for simulation := 0; simulation < NumSimulations; simulation++ {
			batterySimulation := simulator.SimulateBatteryVoltageWithFailure(14, NumSamples)
			dailySimulation := simulator.SimulateDailyTemperature(25, NumSamples)
			heatingSimulation := simulator.SimulateTemperatureWithControl(25, NumSamples, "heating", 30)
			coolingSimulation := simulator.SimulateTemperatureWithControl(25, NumSamples, "cooling", -15)

			var samples []Sample
			for i := 0; i < NumSamples; i++ {
				sample := Sample{
					StationId: stationId,
					Timestamp: date,
					Voltage:   batterySimulation[i],
					Daily:     dailySimulation[i],
					Heating:   heatingSimulation[i],
					Cooling:   coolingSimulation[i],
				}
				samples = append(samples, sample)
				date = date.Add(-10 * time.Minute)
			}
			simulations = append(simulations, samples)
		}
		result = append(result, simulations)
	}
	return result
}

const (
	influxDBToken  = "zoAUTXUsH2U1BrQgQwU8OQmF__b9g6UqFnxirpVLWAauGORjeiBVMsiuDFRFiNa53PVGNXa5Ia4L3DDvbEm0_w=="
	influxDBURL    = "http://localhost:8086"
	influxDBOrg    = "zbyju"
	influxDBBucket = "benchmark"
)

func writeToInfluxDB(data [][][]Sample) {
	// Create a new client using an InfluxDB server base URL and an authentication token
	client := influxdb2.NewClient(influxDBURL, influxDBToken)
	defer client.Close()

	// Get non-blocking write client
	writeAPI := client.WriteAPIBlocking(influxDBOrg, influxDBBucket)

	for stationId, stationData := range data {
		for simId, simulation := range stationData {
			fmt.Println("Done simulation ", simId, ", station ", stationId, " data: ", simulation[0])
			for _, sample := range simulation {
				// Create a new point using full params constructor
				p := influxdb2.NewPointWithMeasurement("environment_data").
					AddTag("stationId", fmt.Sprintf("%d", sample.StationId)).
					AddField("voltage", sample.Voltage).
					AddField("daily", sample.Daily).
					AddField("heating", sample.Heating).
					AddField("cooling", sample.Cooling).
					SetTime(sample.Timestamp)

				// Write the point immediately
				err := writeAPI.WritePoint(context.Background(), p)
				if err != nil {
					fmt.Printf("Write error: %s\n", err.Error())
				}
			}
		}
	}
}

const (
	tshost     = "localhost" // or Docker network address if applicable
	tsport     = 5432
	tsuser     = "postgres"
	tspassword = "password"
)

func setupTimescaleDB() {
	// Connection string to the default PostgreSQL database
	defaultDBConnStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable", tshost, tsport, tsuser, tspassword)

	// Connect to the default PostgreSQL database
	conn, err := pgx.Connect(context.Background(), defaultDBConnStr)
	if err != nil {
		log.Fatalf("Unable to connect to default database: %v", err)
	}
	defer conn.Close(context.Background())

	// Create the benchmark database (if it doesn't exist)
	_, err = conn.Exec(context.Background(), "CREATE DATABASE benchmark;")
	if err != nil {
		log.Printf("Database creation might have failed (it's okay if it already exists): %v", err)
	}

	// Connection string to the benchmark database
	benchmarkDBConnStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=benchmark sslmode=disable", tshost, tsport, tsuser, tspassword)

	// Connect to the benchmark database
	benchmarkConn, err := pgx.Connect(context.Background(), benchmarkDBConnStr)
	if err != nil {
		log.Fatalf("Unable to connect to benchmark database: %v", err)
	}
	defer benchmarkConn.Close(context.Background())

	// Create environment_data table and setup indexes
	setupSQL := `
	CREATE TABLE IF NOT EXISTS environment_data (
		station_id INT,
		timestamp TIMESTAMPTZ NOT NULL,
		voltage DOUBLE PRECISION,
		daily DOUBLE PRECISION,
		cooling DOUBLE PRECISION,
		heating DOUBLE PRECISION
	);

	CREATE INDEX IF NOT EXISTS idx_timestamp ON environment_data (timestamp);
	CREATE INDEX IF NOT EXISTS idx_station_id ON environment_data (station_id);

	`
	_, err = benchmarkConn.Exec(context.Background(), setupSQL)
	if err != nil {
		log.Fatalf("Failed to create table, indexes, or convert to hypertable: %v", err)
	}
}

const (
	pshost     = "localhost" // or Docker network address if applicable
	psport     = 5432
	psuser     = "postgres"
	pspassword = "password"
)

func setupPostgreSQL() {
	// Connection string to the default PostgreSQL database
	defaultDBConnStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable", pshost, psport, psuser, pspassword)

	// Connect to the default PostgreSQL database
	conn, err := pgx.Connect(context.Background(), defaultDBConnStr)
	if err != nil {
		log.Fatalf("Unable to connect to default database: %v", err)
	}
	defer conn.Close(context.Background())

	// Create the benchmark database (if it doesn't exist)
	_, err = conn.Exec(context.Background(), "CREATE DATABASE benchmark;")
	if err != nil {
		log.Printf("Database creation might have failed (it's okay if it already exists): %v", err)
	}

	// Connection string to the benchmark database
	benchmarkDBConnStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=benchmark sslmode=disable", pshost, psport, psuser, pspassword)

	// Connect to the benchmark database
	benchmarkConn, err := pgx.Connect(context.Background(), benchmarkDBConnStr)
	if err != nil {
		log.Fatalf("Unable to connect to benchmark database: %v", err)
	}
	defer benchmarkConn.Close(context.Background())

	// Create environment_data table and setup indexes
	setupSQL := `
	CREATE TABLE IF NOT EXISTS environment_data (
		station_id INT,
		timestamp TIMESTAMPTZ NOT NULL,
		voltage DOUBLE PRECISION,
		daily DOUBLE PRECISION,
		cooling DOUBLE PRECISION,
		heating DOUBLE PRECISION
	);

	CREATE INDEX IF NOT EXISTS idx_timestamp ON environment_data (timestamp);
	CREATE INDEX IF NOT EXISTS idx_station_id ON environment_data (station_id);
	`
	_, err = benchmarkConn.Exec(context.Background(), setupSQL)
	if err != nil {
		log.Fatalf("Failed to create table, indexes, or convert to hypertable: %v", err)
	}
}

const (
	tsConnString = "postgres://postgres:password@localhost:5432/benchmark"
	pgConnString = "postgres://postgres:password@localhost:5433/benchmark"
)

func writeToTimescaleDB(data [][][]Sample, connString string) {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		return
	}
	defer conn.Close(context.Background())

	for _, stationData := range data {
		for _, simulation := range stationData {
			for _, sample := range simulation {
				_, err := conn.Exec(context.Background(), "INSERT INTO environment_data (station_id, timestamp, voltage, daily, heating, cooling) VALUES ($1, $2, $3, $4, $5, $6)",
					sample.StationId, sample.Timestamp, sample.Voltage, sample.Daily, sample.Heating, sample.Cooling)
				if err != nil {
					fmt.Printf("Insert error: %v\n", err)
				}
			}
		}
	}
}

/* func BenchmarkAddInfluxDB(b *testing.B) {
	// setup code here
	data := generateData()

	for i := 0; i < b.N; i++ {
		// perform the operation you're benchmarking
		writeToInfluxDB(data)

	}

	// teardown code here
}
*/

func BenchmarkAddTimescaleDB(b *testing.B) {
	// setup code here
	setupTimescaleDB()
	data := generateData()

	for i := 0; i < b.N; i++ {
		// perform the operation you're benchmarking
		writeToTimescaleDB(data, tsConnString)

	}

	// teardown code here
}
