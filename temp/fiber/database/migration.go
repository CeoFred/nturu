package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nturu/microservice-template/internal/models"
	"gorm.io/gorm"
)

const (
	dbTimeout = 30 * time.Second
)

func RunManualMigration(db *gorm.DB) {

	query1 := `CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			user_id VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL,
			first_name VARCHAR(255) NOT NULL,
			last_name VARCHAR(255) NOT NULL,
			ip VARCHAR(255)DEFAULT NULL,
			account_status INTEGER DEFAULT 1 NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			last_login VARCHAR(255) NULL,
		 	account_type VARCHAR(255) NULL,
		 role VARCHAR(255) DEFAULT 'USER',
		 email_verified BOOLEAN DEFAULT FALSE,
		 country VARCHAR(255) DEFAULT NULL,
		 phone_number VARCHAR(255) DEFAULT NULL,
		 status VARCHAR(255) DEFAULT 'Inactive'
			);`

	migrationQueries := []string{
		query1,
	}

	log.Println("running db migration :::::::::::::")

	_, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	for _, query := range migrationQueries {
		err := db.Exec(query).Error
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("complete db migration")
}

// this should handle gorm auto migration
func RunAutoMigrations() {
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		panic(fmt.Errorf("failed to migrate: %s", err))
	}
	fmt.Println("Migrations completed")
}
