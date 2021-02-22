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
		"http://www.readcomics.tv/the-unbeatable-squirrel-girl/chapter-3/24":    "www.readcomics.tv/the-unbeatable-squirrel-girl",
		"http://readcomiconline.to/Comic/Jughead-2015/Issue-15?id=112397":       "readcomiconline.to/Comic/Jughead-2015",
		"https://kimcartoon.to/Cartoon/DuckTales-2017/Episode-1?id=76477":       "kimcartoon.to/Cartoon/DuckTales-2017",
		"http://pixa.club/en/gravity-falls/season-1/epizod-14-bottomless-pit":   "pixa.club/en/gravity-falls",
		"https://www.youtube.com/watch?v=FkSFBcjOKHY&list=RDFkSFBcjOKHY":        "youtube.com/watch?list=RDFkSFBcjOKHY",
		"https://m.youtube.com/watch?v=FkSFBcjOKHY&list=RDFkSFBcjOKHY":          "youtube.com/watch?list=RDFkSFBcjOKHY",
		"https://youtu.be/watch?v=FkSFBcjOKHY&list=RDFkSFBcjOKHY":               "youtube.com/watch?list=RDFkSFBcjOKHY",
		"https://youtu.be/watch?v=FkSFBcjOKHY":                                  "youtube.com/watch?v=FkSFBcjOKHY",
		"https://www.facebook.com/humandiscoveriesshow/videos/501578117278841/": "facebook.com/humandiscoveriesshow",
		"https://www.gocomics.com/sarahs-scribbles/2014/01/02":                  "www.gocomics.com/sarahs-scribbles",
		"http://devdocs.io/rust/book/first-edition/primitive-types":             "devdocs.io/rust",
	} {

		u, _ := url.Parse(uri)

		got := urlToUniqKey(u)

		if got != expected {
			t.Error("Expected ", expected, "Got ", got)
		}
	}
}
