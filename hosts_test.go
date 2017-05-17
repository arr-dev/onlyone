package main

import (
	"net/url"
	"testing"
)

func TestSimpleHostKey(t *testing.T) {
	u, _ := url.Parse("http://oglaf.com/apocrypha/2/")

	expected := "oglaf.com"
	got := urlToUniqKey(u)

	if got != expected {
		t.Error("Expected ", expected, "Got ", got)
	}
}

func TestMultipleHostKey(t *testing.T) {
	for uri, expected := range map[string]string{
		"http://www.readcomics.tv/the-unbeatable-squirrel-girl/chapter-3/24": "www.readcomics.tv/the-unbeatable-squirrel-girl",
		"http://readcomiconline.to/Comic/Jughead-2015/Issue-15?id=112397":    "readcomiconline.to/Comic/Jughead-2015",
	} {

		u, _ := url.Parse(uri)

		got := urlToUniqKey(u)

		if got != expected {
			t.Error("Expected ", expected, "Got ", got)
		}
	}
}
