package postgresql

// PostgreSQL

func InitDatabasePostgreSQL() string {
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

func InitQueryPostgreSQL() string {
	return `
    CREATE TABLE IF NOT EXISTS iot_snapshots (
        station_id INT NOT NULL,
        created_at TIMESTAMP NOT NULL,
        outside_temperature DECIMAL,
        voltage DECIMAL,
        heating_temperature DECIMAL,
        cooling_temperature DECIMAL
    );

    CREATE INDEX IF NOT EXISTS idx_stationid ON iot_snapshots (station_id);
    CREATE INDEX IF NOT EXISTS idx_createdat ON iot_snapshots (created_at);
  `
}

func InsertQueryPostgreSQL() string {
	return `
    INSERT INTO iot_snapshots (station_id, created_at, outside_temperature, voltage, heating_temperature, cooling_temperature)
    VALUES ($1, $2, $3, $4, $5, $6)
  `
}

func GetQueryPostgreSQL() string {
	return `
    SELECT * FROM iot_snapshots
    WHERE station_id = 1 AND created_at >= NOW() - INTERVAL '1 day';
  `
}

func AvgQueryPostgreSQL() string {
	return `
    SELECT 
      AVG(outside_temperature) AS avg_outside_temperature,
      AVG(voltage) AS avgVoltage,
      AVG(heating_temperature) AS avg_heating_temperature,
      AVG(cooling_temperature) AS avg_cooling_temperature
    FROM iot_snapshots
    WHERE station_id = 1;
  `
}
