package main

import (
	"database/sql"
	"errors"
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
	err := checkAuth(w, r)
	if err != nil {
		log.Printf("auth failed %q", err)
		return
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	log.Println("db")
	handleErr(err)
	err = db.Ping()
	log.Println("db1")
	handleErr(err)
	log.Println("db2")

	switch r.Method {
	case "GET":
		list(db, w)
	case "POST":
		linkId := r.FormValue("id")
		toDelete := r.FormValue("delete")
		link := r.FormValue("link")

		if linkId != "" && toDelete == "true" {
			remove(linkId, db)
		}

		u, err := url.Parse(link)
		if link != "" && err == nil {
			add(u, db)
		}

		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
}

func list(db *sql.DB, w http.ResponseWriter) {
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
}

func remove(linkId string, db *sql.DB) {
	log.Printf("remove link %d", linkId)
	stmt, err := db.Prepare("DELETE FROM links WHERE id = $1")
	handleErr(err)

	_, err = stmt.Exec(linkId)
	handleErr(err)
}

func add(u *url.URL, db *sql.DB) {
	hostId := getHost(u.Host, db)

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

func getHost(host string, db *sql.DB) int {
	log.Println("host: " + host)
	var hostId int
	err := db.QueryRow("SELECT id FROM hosts WHERE host = $1", host).Scan(&hostId)

	if err == sql.ErrNoRows {
		log.Println("create host")
		err = db.QueryRow("INSERT INTO hosts (host, created_at, updated_at) values($1, $2, $2) RETURNING id", host, time.Now().UTC()).Scan(&hostId)
		handleErr(err)
	}
	handleErr(err)
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
