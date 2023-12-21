package main

import (
	"benchmark/logic/influxdb"
	"benchmark/logic/mongodb"
	"benchmark/logic/postgresql"
	"benchmark/logic/questdb"
	"benchmark/logic/timescaledb"
	"benchmark/simulator"
	"context"
	"fmt"
	"testing"
)

const (
	numSnapshots   = 144
	numSimulations = 200
	numStations    = 100
)

func BenchmarkAddInfluxDB(b *testing.B) {
	// setup code here
	client, writeAPI, _ := influxdb.Init()
	data := simulator.GenerateData(numSnapshots, numSimulations, numStations)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// perform the operation you're benchmarking
		err := influxdb.InsertAll(writeAPI, data)
		if err != nil {
			fmt.Println("Error when inserting data: ", err)
		}
	}
	b.StopTimer()

	// teardown code here
	fmt.Println(b.Name() + " took: " + fmt.Sprint(b.Elapsed().Abs()))
	client.Close()
}

func BenchmarkGetFlux(b *testing.B) {
	// setup code here
	client, _, queryAPI := influxdb.Init()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// perform the operation you're benchmarking
		_, err := influxdb.QueryData(queryAPI, influxdb.GetQueryFlux())
		if err != nil {
			fmt.Println("Error when getting data: ", err)
		}
	}
	b.StopTimer()

	// teardown code here
	fmt.Println(b.Name() + " took: " + fmt.Sprint(b.Elapsed().Abs()))
	client.Close()
}

func BenchmarkAvgFlux(b *testing.B) {
	// setup code here
	client, _, queryAPI := influxdb.Init()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// perform the operation you're benchmarking
		_, err := influxdb.QueryData(queryAPI, influxdb.AggregateQueryFlux())
		if err != nil {
			fmt.Println("Error when getting data: ", err)
		}
	}
	b.StopTimer()

	// teardown code here
	fmt.Println(b.Name() + " took: " + fmt.Sprint(b.Elapsed().Abs()))
	client.Close()
}

func BenchmarkAddPostgreSQL(b *testing.B) {
	// setup code here
	pool := postgresql.Init()
	data := simulator.GenerateData(numSnapshots, numSimulations, numStations)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// perform the operation you're benchmarking
		err := postgresql.InsertAll(pool, data)
		if err != nil {
			fmt.Println("Error when inserting data: ", err)
		}
	}
	b.StopTimer()

	// teardown code here
	fmt.Println(b.Name() + " took: " + fmt.Sprint(b.Elapsed().Abs()))
}

func BenchmarkGetPostgreSQL(b *testing.B) {
	// setup code here
	pool := postgresql.Init()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// perform the operation you're benchmarking
		_, err := postgresql.QueryData(pool, postgresql.GetQueryPostgreSQL())
		if err != nil {
			fmt.Println("Error when getting data: ", err)
		}
	}
	b.StopTimer()

	// teardown code here
	fmt.Println(b.Name() + " took: " + fmt.Sprint(b.Elapsed().Abs()))
}

func BenchmarkAvgPostgreSQL(b *testing.B) {
	// setup code here
	pool := postgresql.Init()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// perform the operation you're benchmarking
		_, err := postgresql.QueryData(pool, postgresql.AvgQueryPostgreSQL())
		if err != nil {
			fmt.Println("Error when getting data: ", err)
		}
	}
	b.StopTimer()

	// teardown code here
	fmt.Println(b.Name() + " took: " + fmt.Sprint(b.Elapsed().Abs()))
}

func BenchmarkAddTimescaleDB(b *testing.B) {
	// setup code here
	pool := timescaledb.Init()
	data := simulator.GenerateData(numSnapshots, numSimulations, numStations)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// perform the operation you're benchmarking
		err := timescaledb.InsertAll(pool, data)
		if err != nil {
			fmt.Println("Error when inserting data: ", err)
		}
	}
	b.StopTimer()

	// teardown code here
	fmt.Println(b.Name() + " took: " + fmt.Sprint(b.Elapsed().Abs()))
}

func BenchmarkGetTimescaleDB(b *testing.B) {
	// setup code here
	pool := timescaledb.Init()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// perform the operation you're benchmarking
		_, err := timescaledb.QueryData(pool, timescaledb.GetQueryTimescaleDB())
		if err != nil {
			fmt.Println("Error when getting data: ", err)
		}
	}
	b.StopTimer()

	// teardown code here
	fmt.Println(b.Name() + " took: " + fmt.Sprint(b.Elapsed().Abs()))
}

func BenchmarkAvgTimescaleDB(b *testing.B) {
	// setup code here
	pool := timescaledb.Init()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// perform the operation you're benchmarking
		_, err := timescaledb.QueryData(pool, timescaledb.AvgQueryTimescaleDB())
		if err != nil {
			fmt.Println("Error when getting data: ", err)
		}
	}
	b.StopTimer()

	// teardown code here
	fmt.Println(b.Name() + " took: " + fmt.Sprint(b.Elapsed().Abs()))
}

func BenchmarkAddQuestDB(b *testing.B) {
	// setup code here
	sender, ctx, pool := questdb.Init()
	data := simulator.GenerateData(numSnapshots, numSimulations, numStations)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// perform the operation you're benchmarking
		err := questdb.InsertAll(sender, ctx, data)
		if err != nil {
			fmt.Println("Error when inserting data: ", err)
		}
	}
	b.StopTimer()

	// teardown code here
	fmt.Println(b.Name() + " took: " + fmt.Sprint(b.Elapsed().Abs()))
	sender.Close()
	pool.Close()
}

func BenchmarkGetQuestDB(b *testing.B) {
	// setup code here
	_, _, pool := questdb.Init()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// perform the operation you're benchmarking
		_, err := questdb.QueryData(pool, questdb.GetQueryQuestDB())
		if err != nil {
			fmt.Println("Error when getting data: ", err)
		}
	}
	b.StopTimer()

	// teardown code here
	fmt.Println(b.Name() + " took: " + fmt.Sprint(b.Elapsed().Abs()))
}

func BenchmarkAvgQuestDB(b *testing.B) {
	// setup code here
	_, _, pool := questdb.Init()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// perform the operation you're benchmarking
		_, err := questdb.QueryData(pool, questdb.AvgQueryQuestDB())
		if err != nil {
			fmt.Println("Error when getting data: ", err)
		}
	}
	b.StopTimer()

	// teardown code here
	fmt.Println(b.Name() + " took: " + fmt.Sprint(b.Elapsed().Abs()))
}

func BenchmarkAddMongoDB(b *testing.B) {
	// setup code here
	client := mongodb.Init()
	data := simulator.GenerateData(numSnapshots, numSimulations, numStations)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// perform the operation you're benchmarking
		err := mongodb.InsertAll(client, data)
		if err != nil {
			fmt.Println("Error when inserting data: ", err)
		}
	}
	b.StopTimer()

	// teardown code here
	fmt.Println(b.Name() + " took: " + fmt.Sprint(b.Elapsed().Abs()))
	client.Disconnect(context.Background())
}

func BenchmarkGetMongoDB(b *testing.B) {
	// setup code here
	client := mongodb.Init()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// perform the operation you're benchmarking
		_, err := mongodb.QueryLastDayData(client)
		if err != nil {
			fmt.Println("Error when inserting data: ", err)
		}
	}
	b.StopTimer()

	// teardown code here
	fmt.Println(b.Name() + " took: " + fmt.Sprint(b.Elapsed().Abs()))
	client.Disconnect(context.Background())
}

func BenchmarkAvgMongoDB(b *testing.B) {
	// setup code here
	client := mongodb.Init()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// perform the operation you're benchmarking
		_, err := mongodb.QueryAverageValues(client)
		if err != nil {
			fmt.Println("Error when inserting data: ", err)
		}
	}
	b.StopTimer()

	// teardown code here
	fmt.Println(b.Name() + " took: " + fmt.Sprint(b.Elapsed().Abs()))
	client.Disconnect(context.Background())
}
