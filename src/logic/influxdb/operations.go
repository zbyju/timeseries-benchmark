package influxdb

import (
	"benchmark/logic"
	"context"
	"fmt"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

const (
	influxToken = "token"
	influxURL   = "http://localhost:8086"
	org         = "org"
	bucket      = "benchmark"
)

func Init() (influxdb2.Client, api.WriteAPI, api.QueryAPI) {
	client := influxdb2.NewClient(influxURL, influxToken)
	writeAPI := client.WriteAPI(org, bucket)
	queryAPI := client.QueryAPI(org)
	return client, writeAPI, queryAPI
}

func Insert(writeAPI api.WriteAPI, snapshot logic.Snapshot) {
	p := influxdb2.NewPointWithMeasurement("iot_snapshots").
		AddTag("stationId", fmt.Sprintf("%d", snapshot.StationID)).
		AddField("outsideTemperature", snapshot.OutsideTemperature).
		AddField("voltage", snapshot.Voltage).
		AddField("heatingTemperature", snapshot.HeatingTemperature).
		AddField("coolingTemperature", snapshot.CoolingTemperature).
		SetTime(snapshot.CreatedAt)

	writeAPI.WritePoint(p)
}

func InsertAll(writeAPI api.WriteAPI, data [][][]logic.Snapshot) error {
	for _, station := range data {
		for _, simulation := range station {
			for _, snapshot := range simulation {
				Insert(writeAPI, snapshot)
			}
		}
	}
	writeAPI.Flush()
	return nil
}

func QueryData(queryAPI api.QueryAPI, query string) (*api.QueryTableResult, error) {
	return queryAPI.Query(context.Background(), query)
}
