package timescaledb

func InitDatabaseTimescaleDB() string {
	return `
  DO
  $do$
  DECLARE
    _db TEXT := 'benchmark';
    _user TEXT := 'postgres';
    _password TEXT := 'password';
  BEGIN
    CREATE EXTENSION IF NOT EXISTS dblink; -- enable extension 
    IF EXISTS (SELECT 1 FROM pg_database WHERE datname = _db) THEN
      RAISE NOTICE 'Database already exists';
    ELSE
      PERFORM dblink_connect('host=localhost user=' || _user || ' password=' || _password || ' dbname=' || current_database());
      PERFORM dblink_exec('CREATE DATABASE ' || _db);
    END IF;
  END
  $do$
  `
}

func InitQueryTimescaleDB() string {
	return `
    DROP TABLE IF EXISTS iot_snapshots CASCADE;
    CREATE TABLE IF NOT EXISTS iot_snapshots (
        station_id INT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL,
        outside_temperature DECIMAL,
        voltage DECIMAL,
        heating_temperature DECIMAL,
        cooling_temperature DECIMAL
    );
    SELECT create_hypertable('iot_snapshots', 'created_at', create_default_indexes => FALSE);
    CREATE INDEX IF NOT EXISTS idx_stationid ON iot_snapshots (station_id);
    CREATE INDEX IF NOT EXISTS idx_createdat ON iot_snapshots (created_at);
    `
}

func CleanQueryTimescaleDB() string {
	return `
    DROP TABLE IF EXISTS iot_snapshots CASCADE;
    `
}

func InsertQueryTimescaleDB() string {
	return `
    INSERT INTO iot_snapshots (station_id, created_at, outside_temperature, voltage, heating_temperature, cooling_temperature)
    VALUES ($1, $2, $3, $4, $5, $6);
    `
}

func GetQueryTimescaleDB() string {
	return `
    SELECT * FROM iot_snapshots
    WHERE station_id = 1 AND created_at >= NOW() - INTERVAL '1 day';
    `
}

func AvgQueryTimescaleDB() string {
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
