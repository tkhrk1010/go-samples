package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
)

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func main() {
	albums, err := albumsByArtist("John Coltrane")
	if err != nil {
			log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albums)
}

func connectDB() (*sql.DB, error){
    // Capture connection properties.
		connStr := "user=sample dbname=recordings sslmode=disable"

    // Get a database handle.
    var err error
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }

    pingErr := db.Ping()
    if pingErr != nil {
        log.Fatal(pingErr)
    }
    fmt.Println("Connected!")

		return db, nil
}

// albumsByArtist queries for albums that have the specified artist name.
func albumsByArtist(name string) ([]Album, error) {
	// An albums slice to hold data from returned rows.
	var albums []Album

	db, err := connectDB()
	if err != nil {
		return nil, fmt.Errorf("DB connection error: %v", err)
  }

	rows, err := db.Query("SELECT * FROM album WHERE artist = $1", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}

	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
			var alb Album
			if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
					return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
			}
			albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}