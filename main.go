package main

import (
	d "longevity/src/database"
	. "longevity/src/rest"
)

func main() {
	db := &d.DB{Path: "longevity.db"}

	sql_db := d.Run(db)
	defer sql_db.Close()

	var rest API
	rest = NewRestInterface(sql_db)
	rest.Run()
}
