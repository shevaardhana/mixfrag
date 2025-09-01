package database

import (
	"fmt"
	"log"
)

func Migrate() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(100) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			created_by VARCHAR(100),
			modified_at TIMESTAMP,
			modified_by VARCHAR(100)
		);`,
		`INSERT INTO users (username, password, created_at, created_by)
			VALUES ('sheva', 'sheva123', NOW(), 'system')
			ON CONFLICT (username) DO NOTHING;`,
	}

	for _, q := range queries {
		_, err := DB.Exec(q)
		if err != nil {
			log.Fatal("Failed to migrate table:", err)
		}
	}

	fmt.Println("✅ Database migration completed")
}
