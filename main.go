package main

import (
	. "longevity/src/database"
	. "longevity/src/rest"
)

func main() {
	db := SetupPostgresDB("postgres", "foobar", "postgres")
	// db := SetupSQLiteDB("longevity.db", "longevity")

	defer db.Close()

	sample := NewDevice("anker", "01:23:45:67:89", "es-lite", "0.0.1")
	sample.CreateDevice(db)

	var rest API
	rest = NewRestInterface(db)
	go rest.Run(8000)

	var rest2 API
	rest2 = NewRestInterface(db)
	rest2.RunWithHTTPS(443)
}
