package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/url"
	"os"
)

const FileName = "file.json"

func main() {

	flag.Parse()
	uri := flag.Arg(0)

	u, _ := url.Parse(uri)

	db, _ := readDB()

	if db == nil {
		db = make(map[string]string)
	}

	db[u.Host] = uri

	writeDB(db)
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
