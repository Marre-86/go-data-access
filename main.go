package main

import (
	"fmt"
	"os"
	"log"
	"bufio"
	"strings"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

type Album struct {
	ID	int64
	Title	string
	Artist	string
	Price	float32
}

func main() {
	connString := "postgres://marre:Hombre1986@127.0.0.1:5432/recordings"
	var err error
	db, err = sql.Open("pgx", connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Connected!")
	defer db.Close()


	fmt.Println("Choose type of search")
	fmt.Println("1. Search by artist name")
	fmt.Println("2. Search by album ID")
	fmt.Println("3. Add album to a database")
	fmt.Println("Enter your choice (1 or 2 or 3)")

	var inp int64
	_, err = fmt.Scanln(&inp)
	handleInputError(err)

	if inp == 1 {
		searchByArtist()
	} else if inp == 2 {
		searchByAlbumID()
	} else if inp == 3 {
		insertAlbum()
	} else {
		fmt.Println("Invalid choice")
	}
}

func searchByArtist() {
	fmt.Println("Enter the name of the artist")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	handleInputError(err)
	input = strings.TrimSpace(input)
	albums, err := albumsByArtist(input)
	fmt.Printf("Looking for albums of arist: %v\n", input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albums)
}

func searchByAlbumID() {
	fmt.Println("Enter the album ID")
	var inp int64
	_, err := fmt.Scanln(&inp)
	handleInputError(err)
	album, err := albumByID(inp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found: %v\n", album)
}

func insertAlbum() {
	fmt.Println("Enter the name of the artist")
	reader := bufio.NewReader(os.Stdin)
	inpArtist, err := reader.ReadString('\n')
	handleInputError(err)
	inpArtist = strings.TrimSpace(inpArtist)
	fmt.Println("Enter the name of the album")
	reader = bufio.NewReader(os.Stdin)
	inpTitle, err := reader.ReadString('\n')
	handleInputError(err)
	inpTitle = strings.TrimSpace(inpTitle)
	var inpPrice float32
	fmt.Println("Enter the price")
	_, err = fmt.Scanf("%f", &inpPrice)
	handleInputError(err)

	albID, err := addAlbum(Album{
		Title:	inpTitle,
		Artist:	inpArtist,
		Price: inpPrice,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Success! ID of added album: %v\n", albID)
}


func albumsByArtist(name string) ([]Album, error) {
	var albums []Album
	rows, err := db.Query("SELECT * FROM album WHERE artist = $1", name)
	if  err != nil {
		return nil, fmt.Errorf("111albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err!= nil {
			return  nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}

	return albums, nil
}

func albumByID(id int64) (Album, error) {
	var alb Album
	row := db.QueryRow("SELECT * FROM album WHERE id = $1", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumById %d: no such album", id)
		}
		return alb, fmt.Errorf("albumById %d: no such album", id)
	}
	return alb, nil
}

func addAlbum(alb Album) (int64, error) {
	var id int64
	err := db.QueryRow("INSERT INTO album (title, artist, price) VALUES ($1, $2, $3) RETURNING id", alb.Title, alb.Artist, alb.Price).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}

func handleInputError(err error) {
	if err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}
}

