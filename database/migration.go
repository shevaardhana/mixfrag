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

		`CREATE TABLE category_notes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    desc TEXT,
    created_at TIMESTAMP NOT NULL,
    created_by VARCHAR(255) NOT NULL,
    modified_at TIMESTAMP NULL,
    modified_by VARCHAR(255) NULL
		);`,

		`CREATE TABLE category_parfume (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    desc TEXT,
    created_at TIMESTAMP NOT NULL,
    created_by VARCHAR(255) NOT NULL,
    modified_at TIMESTAMP NULL,
    modified_by VARCHAR(255) NULL
		);`,

		`CREATE TABLE category_parfumes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    desc TEXT,
    created_at TIMESTAMP NOT NULL,
    created_by VARCHAR(255) NOT NULL,
    modified_at TIMESTAMP NULL,
    modified_by VARCHAR(255) NULL
		);`,

		`CREATE TABLE category_smells (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    desc TEXT,
    created_at TIMESTAMP NOT NULL,
    created_by VARCHAR(255) NOT NULL,
    modified_at TIMESTAMP NULL,
    modified_by VARCHAR(255) NULL
		);`,

		`CREATE TABLE essential_oils (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category_smell_id INT NOT NULL REFERENCES category_smells(id),
    category_note_id INT NOT NULL REFERENCES category_notes(id),
    rank_note VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    created_by VARCHAR(255) NOT NULL,
    modified_at TIMESTAMP NULL,
    modified_by VARCHAR(255) NULL
		);`,

		`CREATE TABLE parfumes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    total_ml VARCHAR(50) NOT NULL,
    category_parfume_id INT NOT NULL REFERENCES category_parfumes(id),
    total_oil_drop INT NOT NULL,
    total_oil INT NOT NULL,
    total_parfume_base INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    created_by VARCHAR(255) NOT NULL,
    modified_at TIMESTAMP NULL,
    modified_by VARCHAR(255) NULL
		);`,

		`CREATE TABLE parfume_details (
    id SERIAL PRIMARY KEY,
    parfume_id INT NOT NULL REFERENCES parfumes(id) ON DELETE CASCADE,
    oil_id VARCHAR(100) NOT NULL,
    total_drop VARCHAR(50) NOT NULL,
    rasio DECIMAL(10,4) NOT NULL,
    weight_total INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    created_by VARCHAR(255) NOT NULL,
    modified_at TIMESTAMP NULL,
    modified_by VARCHAR(255) NULL
		);`,
	}

	for _, q := range queries {
		_, err := DB.Exec(q)
		if err != nil {
			log.Fatal("Failed to migrate table:", err)
		}
	}

	fmt.Println("✅ Database migration completed")
}
