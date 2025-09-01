package main

import (
	"mixfrag/database"
	"mixfrag/routers"
)

func main() {
	database.InitDB()
	database.Migrate()
	r := routers.SetupRouter()
	r.Run(":8080")
}
