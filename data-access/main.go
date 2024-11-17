package main

import (
	"database/sql"
	"fmt"
	"log"
	//"os"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Album struct {
	ID int64
	Title string
	Artist string
	Price float32
}

func main() {


	cfg := mysql.Config{
		User:    "root",
        Passwd:  "helsasp",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings", 
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	albums, err := albumsByArtist("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albums)

	// Hard-code ID 2 to test
	alb,err := albumByID(2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found : %v\n",alb)

	// add new album
	albID, err := addAlbum(Album {
		Title: "The Modern Sound of Betty Carter",
		Artist : "Betty Carter",
		Price : 48.99,
	})
	if err != nil {
		log.Fatal(err) 
	}
	fmt.Printf("ID of added album : %v\n", albID)
}

// queries for albums with specicified artist
func albumsByArtist (name string) ([]Album,error) {
	var albums [] Album

	rows,err := db.Query ("SELECT * FROM album WHERE artist =?",name)
	if err!= nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)

	}
	defer rows.Close()
	// Loop thru rows, use scan to assign column data to struct 
	for rows.Next(){
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

// albumByID queries for the album with the specified ID.
func albumByID(id int64) (Album, error) {
    // An album to hold data from the returned row.
    var alb Album

    row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
    if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
        if err == sql.ErrNoRows {
            return alb, fmt.Errorf("albumsById %d: no such album", id)
        }
        return alb, fmt.Errorf("albumsById %d: %v", id, err)
    }
    return alb, nil
}

// add albums to database, return id of new entry
func addAlbum (alb Album) (int64,error) {
	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("addAlbum : %v", err)
	}
	id,err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id,nil
}