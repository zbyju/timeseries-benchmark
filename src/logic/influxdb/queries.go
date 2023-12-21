package influxdb

func GetQueryFlux() string {
	return `
    from(bucket: "benchmark")
      |> range(start: -1d)
      |> filter(fn: (r) => r._measurement == "iot_snapshots" and r.stationId == "1")
  `
}

func AggregateQueryFlux() string {
	return `
    from(bucket: "benchmark")
      |> range(start: 0)
      |> filter(fn: (r) => r._measurement == "iot_snapshots" and r.stationId == "1")
      |> aggregateWindow(every: inf, fn: mean)
  `
}
