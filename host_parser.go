package main

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
)

type Host struct {
	key string
}

func parseHost(u *url.URL) *Host {
	log.Println("host: " + u.Host)
	var host string

	readcomiconlinePat := regexp.MustCompile(`readcomiconline\..*`)
	kimcartoonPat := regexp.MustCompile(`kimcartoon\..*`)

	switch u.Host {
	case "www.readcomics.tv", "www.gocomics.com", "devdocs.io":
		key := strings.Split(u.Path, "/")[1]
		host = u.Host + "/" + key
	case readcomiconlinePat.FindString(u.Host):
		keys := strings.Split(u.Path, "/")[1:3]
		host = "readcomiconline.to" + "/" + strings.Join(keys, "/")
	case kimcartoonPat.FindString(u.Host):
		keys := strings.Split(u.Path, "/")[1:3]
		host = "kimcartoon.to" + "/" + strings.Join(keys, "/")
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
