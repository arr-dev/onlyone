package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

func fetchIcon(u url.URL) (string, error) {
	log.Printf("fetch icon from %s", u)
	uri, err := attemptFavicon(&u)

	if err != nil {
		uri, err = fetchMetaTag(&u)
	}

	return uri, err
}

func attemptFavicon(u *url.URL) (string, error) {
	u.Path = "/favicon.ico"
	resp, err := http.Head(u.String())
	log.Printf("fetch favicon from %s", u)

	if err == nil && resp.StatusCode == 200 && resp.ContentLength > 0 {
		log.Printf("favicon found %s", u)
		return u.String(), nil
	} else {
		log.Printf("favicon not found %s", resp.Status)
		return "", errors.New(resp.Status)
	}
}

func fetchMetaTag(u *url.URL) (string, error) {
	u.Path = "/"

	log.Printf("fetch meta tag from %s", u)

	resp, err := http.Get(u.String())

	if err == nil {
		var link string
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		pattern := regexp.MustCompile(`<link.*rel="(shortcut )?icon".*>`)
		rel := pattern.FindString(string(body))
		log.Printf("rel %s", rel)
		pattern = regexp.MustCompile(`href="([^" >]+)"`)
		match := pattern.FindStringSubmatch(rel)

		if len(match) > 0 {
			link = match[1]
			log.Printf("link %s", link)
			if link[:1] == "/" {
				u.Path = link
				link = u.String()
			}
		}

		if link != "" {
			return link, nil
		} else {
			err = errors.New("Favicon meta not found")
		}
	}

	log.Printf("failed %q", err)
	return "", err

}
