package main

import "flag"

var (
	httpSpec = flag.String("http", ":8040", "HTTP address and port")
	dbFile   = flag.String("file", "db.json", "Database file")
)
