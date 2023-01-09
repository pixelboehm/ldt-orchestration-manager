package database

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS devices 
(
	id INTEGER PRIMARY KEY, 
	name VARCHAR(64), 
	macAddress VARCHAR(17) UNIQUE, 
	twin VARCHAR(64), 
	version string
)`

const insertDeviceQuery = `INSERT INTO devices 
(
	name, 
	macAddress, 
	twin, 
	version
) VALUES (?, ?, ?, ?)`

const updateDeviceQuery = `UPDATE devices SET 
	name=?, 
	macAddress=?,
	twin=?, 
	version=? 
	WHERE id=?`

const deleteDeviceQuery = `DELETE FROM devices WHERE macAddress=?`

const getDeviceTableQuery = `SELECT name, macAddress, twin, version FROM devices`

const checkIfDeviceExistsQuery = `SELECT EXISTS
(
	SELECT 1 FROM devices WHERE macAddress=?
);`

const getDeviceByIDQuery = `SELECT name, macAddress, twin, version WHERE id = ?`
