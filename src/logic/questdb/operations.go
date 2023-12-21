package questdb

import (
	"benchmark/logic"
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	qdb "github.com/questdb/go-questdb-client/v2"
)

const (
	questDBURL = "postgres://admin:quest@localhost:8812/qdb"
)

func Init() (*qdb.LineSender, context.Context, *pgxpool.Pool) {
	ctx := context.TODO()
	sender, err := qdb.NewLineSender(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize sender: %v\n", err)
		os.Exit(1)
	}

	pool, err := pgxpool.Connect(context.Background(), questDBURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect: %v\n", err)
		os.Exit(1)
	}

	createTableSQL := InitQueryQuestDB()
	_, err = pool.Exec(context.Background(), createTableSQL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize table: %v\n", err)
		os.Exit(1)
	}

	return sender, ctx, pool
}

func Insert(sender *qdb.LineSender, ctx context.Context, snapshot logic.Snapshot) error {
	err := sender.
		Table("iot_snapshots").
		Int64Column("station_id", int64(snapshot.StationID)).
		Float64Column("outside_temperature", snapshot.OutsideTemperature).
		Float64Column("voltage", snapshot.Voltage).
		Float64Column("heating_temperature", snapshot.HeatingTemperature).
		Float64Column("cooling_temperature", snapshot.CoolingTemperature).
		At(ctx, snapshot.CreatedAt)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to add data: %v\n", err)
		return err
	}

	return nil
}

func InsertAll(sender *qdb.LineSender, ctx context.Context, data [][][]logic.Snapshot) error {
	for _, station := range data {
		for _, simulation := range station {
			for _, snapshot := range simulation {
				err := Insert(sender, ctx, snapshot)
				if err != nil {
					return err
				}
			}
		}
	}
	err := sender.Flush(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to flush data: %v\n", err)
		return err
	}
	return nil
}

func QueryData(pool *pgxpool.Pool, sqlQuery string) (pgx.Rows, error) {
	return pool.Query(context.Background(), sqlQuery)
}
