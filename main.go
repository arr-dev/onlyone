package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const FileName = "file.json"

func main() {

	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":54321", nil)

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	db, _ := readDB()
	if db == nil {
		db = make(map[string]string)
	}

	switch r.Method {
	case "GET":
		t, _ := template.ParseFiles("index.html")
		t.Execute(w, db)
	case "POST":
		uri := r.FormValue("url")
		del := r.FormValue("delete")
		u, _ := url.Parse(uri)

		if del == "true" {
			delete(db, u.Host)
		} else {
			db[u.Host] = uri
		}

		writeDB(db)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

}

func readDB() (map[string]string, error) {
	m := make(map[string]string)
	db, err := ioutil.ReadFile(FileName)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(db, &m)

	return m, err
}

func writeDB(db map[string]string) {
	f, _ := os.Create(FileName)
	defer f.Close()

	json, _ := json.Marshal(db)
	f.Write(json)
}
