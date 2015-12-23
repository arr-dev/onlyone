package main

import (
	"database/sql"
	"errors"
	_ "github.com/arr-dev/onlyone/Godeps/_workspace/src/github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Link struct {
	Id       int
	Host     string
	Link     string
	ThumbUrl sql.NullString
}

var db *sql.DB

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	db = connectDb()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/thumb", thumbHandler)
	http.ListenAndServe(":"+port, nil)

}

func handleErr(err error) {
	if err != nil {
		log.Println("got error")
		panic(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := checkAuth(w, r)
	if err != nil {
		log.Printf("auth failed %q", err)
		return
	}

	switch r.Method {
	case "GET":
		list(w)
	case "POST":
		linkId := r.FormValue("id")
		toDelete := r.FormValue("delete")
		link := r.FormValue("link")

		if linkId != "" && toDelete == "true" {
			remove(linkId)
		}

		u, err := url.Parse(link)
		if link != "" && err == nil {
			add(u)
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// PUT /thumb host=explosm.net&uri=http://explosm.net/favicons/favicon.ico
func thumbHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.NotFound(w, r)
		return
	}
	err := checkAuth(w, r)
	if err != nil {
		log.Printf("auth failed %q", err)
		return
	}

	host := r.FormValue("host")
	uri := r.FormValue("uri")

	if host != "" && uri != "" {
		log.Printf("update host %s with thumb uri: %s", host, uri)
		stmt, err := db.Prepare("UPDATE hosts SET thumb_url = $1, updated_at = $2 WHERE host = $3")
		handleErr(err)

		_, err = stmt.Exec(uri, time.Now().UTC(), host)
		handleErr(err)
	}
}

func list(w http.ResponseWriter) {
	data := []Link{}
	rows, err := db.Query(`
		SELECT l.id, h.host, l.link, h.thumb_url FROM links l
		INNER JOIN hosts h ON l.host_id = h.id
	`)
	handleErr(err)
	log.Println("query")
	defer rows.Close()

	log.Println("rows")
	for rows.Next() {
		record := Link{}
		log.Println("scan")
		err = rows.Scan(&record.Id, &record.Host, &record.Link, &record.ThumbUrl)
		handleErr(err)
		data = append(data, record)
	}
	handleErr(rows.Err())
	log.Println("index")
	t, err := template.ParseFiles("index.html")
	handleErr(err)
	t.Execute(w, data)
}

func remove(linkId string) {
	log.Printf("remove link %d", linkId)
	stmt, err := db.Prepare("DELETE FROM links WHERE id = $1")
	handleErr(err)

	_, err = stmt.Exec(linkId)
	handleErr(err)
}

func add(u *url.URL) {
	hostId := getHost(u)

	var linkId int
	err := db.QueryRow("SELECT id FROM links WHERE host_id = $1", hostId).Scan(&linkId)

	if err == sql.ErrNoRows {

		log.Println("create link")
		err := db.QueryRow("INSERT INTO links (host_id, link, created_at, updated_at) VALUES($1, $2, $3, $3) RETURNING id", hostId, u.String(), time.Now().UTC()).Scan(&linkId)
		handleErr(err)
		log.Printf("created link %s", u.String())

	} else if err == nil {

		log.Printf("update link %d", linkId)
		stmt, err := db.Prepare("UPDATE links SET link = $1, updated_at = $2 WHERE id = $3")
		handleErr(err)

		_, err = stmt.Exec(u.String(), time.Now().UTC(), linkId)
		handleErr(err)

	} else {
		handleErr(err)
	}
}

func getHost(u *url.URL) int {
	log.Println("host: " + u.Host)
	var hostId int
	err := db.QueryRow("SELECT id FROM hosts WHERE host = $1", u.Host).Scan(&hostId)

	if err == sql.ErrNoRows {
		log.Println("create host")
		thumb, err := fetchIcon(*u)
		log.Printf("got thumb %s", thumb)
		err = db.QueryRow("INSERT INTO hosts (host, thumb_url, created_at, updated_at) values($1, $2, $3, $3) RETURNING id", u.Host, thumb, time.Now().UTC()).Scan(&hostId)
		handleErr(err)
	} else {
		handleErr(err)
	}
	log.Printf("found host %d", hostId)

	return hostId
}

func checkAuth(w http.ResponseWriter, r *http.Request) error {
	user, pass, ok := r.BasicAuth()

	if ok == true && os.Getenv("AUTH_USER") == user && os.Getenv("AUTH_PASS") == pass {
		log.Printf("auth ok %s", user)
		return nil
	}

	w.Header().Set("WWW-Authenticate", "Basic realm=\"user\"")
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	return errors.New("unauthorized")
}

func connectDb() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	log.Println("db")
	handleErr(err)
	if max, err := strconv.Atoi(os.Getenv("DATABASE_MAX_CONNS")); err == nil && max > 0 {
		log.Printf("max db conns %d", max)
		db.SetMaxOpenConns(max)
	}
	err = db.Ping()
	log.Println("db1")
	handleErr(err)
	log.Println("db2")

	return db
}
