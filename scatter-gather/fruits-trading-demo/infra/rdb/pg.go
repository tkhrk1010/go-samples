package rdb

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
	const (
		host     = "localhost"
		port     = 5433
		user     = "admin"
		password = "scatter_gather_pass" // please use env file if you use this code in production
		dbname   = "scatter_gather_sample"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InsertData(db *sql.DB, tableName string, data map[string]float32) error {
	// Prepare the SQL statement for inserting data
	stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO %s (feature_type, value) VALUES ($1, $2)", tableName))
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Iterate over the map and insert each key-value pair into the table
	for featureType, value := range data {
		_, err := stmt.Exec(featureType, value)
		if err != nil {
			return err
		}
	}

	return nil
}