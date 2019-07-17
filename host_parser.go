package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

type Host struct {
	key string
}

func parseHost(u *url.URL) *Host {
	log.Println("host: " + u.Host)
	var host string

	switch u.Host {
	case "www.readcomics.tv":
		key := strings.Split(u.Path, "/")[1]
		host = u.Host + "/" + key
	case "readcomiconline.to", "kimcartoon.to":
		keys := strings.Split(u.Path, "/")[1:3]
		host = u.Host + "/" + strings.Join(keys, "/")
	case "www.facebook.com", "facebook.com", "m.facebook.com":
		host = "facebook.com"
		keys := strings.Split(u.Path, "/")[1:2]
		host += "/" + strings.Join(keys, "/")
	case "pixa.club":
		keys := strings.Split(u.Path, "/")[1:3]
		host = u.Host + "/" + strings.Join(keys, "/")
	case "www.youtube.com", "youtube.com", "m.youtube.com", "youtu.be":
		host = "youtube.com"
		q := u.Query()

		if list := q.Get("list"); list != "" {
			host += fmt.Sprintf("/watch?list=%s", list)
		} else {
			host += fmt.Sprintf("%s?%s", u.Path, u.RawQuery)
		}
	default:
		host = u.Host
	}

	return &Host{key: host}
}
