package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Link struct {
	Id   int
	Host string
	Link string
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":"+port, nil)

}

func handleErr(err error) {
	if err != nil {
		log.Println("got error")
		panic(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	log.Println("db")
	handleErr(err)
	err = db.Ping()
	log.Println("db1")
	handleErr(err)
	log.Println("db2")

	switch r.Method {
	case "GET":
		data := []Link{}
		rows, err := db.Query(`
		SELECT l.id, h.host, l.link FROM links l
		INNER JOIN hosts h ON l.host_id = h.id
	`)
		handleErr(err)
		log.Println("query")
		defer rows.Close()

		log.Println("rows")
		for rows.Next() {
			record := Link{}
			log.Println("scan")
			err = rows.Scan(&record.Id, &record.Host, &record.Link)
			handleErr(err)
			data = append(data, record)
		}
		handleErr(rows.Err())
		log.Println("index")
		t, err := template.ParseFiles("index.html")
		handleErr(err)
		t.Execute(w, data)
	case "POST":
		link := r.FormValue("link")
		u, err := url.Parse(link)
		if link == "" || err != nil {
			log.Println("empty/invalid link")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		host := u.Host

		log.Println("host: " + host)
		var hostId int
		err = db.QueryRow("SELECT id FROM hosts WHERE host = $1", host).Scan(&hostId)

		if err == sql.ErrNoRows {
			log.Println("create host")
			err = db.QueryRow("INSERT INTO hosts (host, created_at, updated_at) values($1, $2, $2) RETURNING id", host, time.Now().UTC()).Scan(&hostId)
			handleErr(err)
		}
		handleErr(err)
		log.Printf("found host %d", hostId)

		var linkId int
		err = db.QueryRow("INSERT INTO links (host_id, link, created_at, updated_at) VALUES($1, $2, $3, $3) RETURNING id", hostId, link, time.Now().UTC()).Scan(&linkId)
		handleErr(err)
		log.Printf("created link %s", link)

		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

}
