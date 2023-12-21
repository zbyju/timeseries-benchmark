package questdb

func InitQueryQuestDB() string {
	return `
    CREATE TABLE IF NOT EXISTS iot_snapshots (
      station_id SYMBOL,
      created_at TIMESTAMP,
      outside_temperature DOUBLE,
      voltage DOUBLE,
      heating_temperature DOUBLE,
      cooling_temperature DOUBLE
    ), INDEX (station_id CAPACITY 128) TIMESTAMP(created_at);
  `
}

// GetQueryQuestDB creates an SQL query to retrieve the last day's data for a specific station
func GetQueryQuestDB() string {
	return `
    SELECT * FROM iot_snapshots
    WHERE station_id = 1 AND created_at > dateadd('d', -1, now());
  `
}

// AvgQueryQuestDB creates an SQL query to get averages of all fields for a specified station over all time
func AvgQueryQuestDB() string {
	return `
    SELECT 
        AVG(outside_temperature) AS avg_outside_temperature,
        AVG(voltage) AS avg_voltage,
        AVG(heating_temperature) AS avg_heating_temperature,
        AVG(cooling_temperature) AS avg_cooling_temperature
    FROM iot_snapshots
    WHERE station_id = 1;
    `
}

// CleanQueryQuestDB creates an SQL query to delete all data from the measurement
func CleanQueryQuestDB() string {
	return `
    DELETE FROM iot_snapshots;
    `
}
