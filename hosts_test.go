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
	u, _ := url.Parse("http://www.readcomics.tv/the-unbeatable-squirrel-girl/chapter-3/24")

	expected := "www.readcomics.tv/the-unbeatable-squirrel-girl"
	got := urlToUniqKey(u)

	if got != expected {
		t.Error("Expected ", expected, "Got ", got)
	}
}
