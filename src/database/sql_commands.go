package database

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS devices
(
	id SERIAL, 
	name VARCHAR(64), 
	macAddress VARCHAR(17) UNIQUE, 
	twin VARCHAR(64), 
	version VARCHAR(6),
	CONSTRAINT devices_pkey PRIMARY KEY (id)
)`

const insertDeviceQuery = `INSERT INTO devices
(
	name, 
	macAddress, 
	twin, 
	version
) VALUES ($1, $2, $3, $4)`

const updateDeviceQuery = `UPDATE devices SET 
	name=$1, 
	macAddress=$2,
	twin=$3, 
	version=$4 
	WHERE id=$5`

const deleteDeviceQuery = `DELETE FROM devices WHERE id=$1`

const getDeviceTableQuery = `SELECT name, macAddress, twin, version FROM devices`

const checkIfDeviceExistsQuery = `SELECT EXISTS
(
	SELECT 1 FROM devices WHERE macAddress=$1
);`

const getDeviceByIDQuery = `SELECT name, macAddress, twin, version FROM devices WHERE id = $1`
