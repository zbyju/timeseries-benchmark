package timescaledb

import (
	"benchmark/logic"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "benchmark"
)

func Init() *pgxpool.Pool {
	// Connect to the default PostgreSQL database
	defaultPool, err := connectDB("postgres")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to default database: %v\n", err)
		os.Exit(1)
	}
	defer defaultPool.Close()

	// Create the 'benchmark' database
	err = createBenchmarkDB(defaultPool)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create 'benchmark' database: %v\n", err)
		os.Exit(1)
	}

	// Now, connect to the newly created 'benchmark' database
	pool, err := connectDB(dbname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to 'benchmark' database: %v\n", err)
		os.Exit(1)
	}

	// Initialize the database and table
	err = initDatabase(pool)
	if err != nil {
		pool.Close()
		fmt.Fprintf(os.Stderr, "Failed to initialize database: %v\n", err)
		os.Exit(1)
	}

	return pool
}

func createBenchmarkDB(pool *pgxpool.Pool) error {
	// Create the 'benchmark' database
	_, err := pool.Exec(context.Background(), InitDatabaseTimescaleDB())
	return err
}

func connectDB(dbname string) (*pgxpool.Pool, error) {
	// Format the connection string for the 'benchmark' database
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, dbname)
	return pgxpool.Connect(context.Background(), connString)
}

func initDatabase(pool *pgxpool.Pool) error {
	// Create the hypertable
	_, err := pool.Exec(context.Background(), InitQueryTimescaleDB())
	return err
}

func InsertAll(pool *pgxpool.Pool, data [][][]logic.Snapshot) error {
	var valueStrings []string
	var valueArgs []interface{}
	valueIdx := 1

	for _, station := range data {
		for _, simulation := range station {
			for _, snapshot := range simulation {
				valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", valueIdx, valueIdx+1, valueIdx+2, valueIdx+3, valueIdx+4, valueIdx+5))
				valueArgs = append(valueArgs, snapshot.StationID, snapshot.CreatedAt, snapshot.OutsideTemperature, snapshot.Voltage, snapshot.HeatingTemperature, snapshot.CoolingTemperature)
				valueIdx += 6

				if valueIdx > 60000 {
					stmt := fmt.Sprintf("INSERT INTO iot_snapshots (station_id, created_at, outside_temperature, voltage, heating_temperature, cooling_temperature) VALUES %s",
						strings.Join(valueStrings, ","))
					_, err := pool.Exec(context.Background(), stmt, valueArgs...)

					if err != nil {
						return err
					}

					valueStrings = []string{}
					valueArgs = []interface{}{}
					valueIdx = 1
				}
			}
		}
	}

	if len(valueArgs) > 0 {
		stmt := fmt.Sprintf("INSERT INTO iot_snapshots (station_id, created_at, outside_temperature, voltage, heating_temperature, cooling_temperature) VALUES %s",
			strings.Join(valueStrings, ","))
		_, err := pool.Exec(context.Background(), stmt, valueArgs...)
		return err
	}
	return nil
}

func Clean(pool *pgxpool.Pool) {
	query := CleanQueryTimescaleDB()
	_, err := pool.Exec(context.Background(), query)

	if err != nil {
		pool.Close()
		fmt.Fprintf(os.Stderr, "Failed to clean database: %v\n", err)
		os.Exit(1)
	}
}

func QueryData(pool *pgxpool.Pool, sqlQuery string) (pgx.Rows, error) {
	return pool.Query(context.Background(), sqlQuery)
}
