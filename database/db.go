package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "127.0.0.1"
	port     = 5433
	user     = "postgres"
	password = "admin1"
	dbname   = "db_fragrance"
)

var DB *sql.DB

func InitDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Gagal buka koneksi:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Gagal ping DB:", err)
	}

	fmt.Println("✅ Successfully connected to database")
	return DB
}
