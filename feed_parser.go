package main

import (
	"encoding/xml"
	"regexp"

	"github.com/mmcdole/gofeed"
)

// RSS structure to generate filtered RSS feed
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	DC      string   `xml:"xmlns:dc,attr"`
	Content string   `xml:"xmlns:content,attr"`
	Atom    string   `xml:"xmlns:atom,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	AtomLink      AtomLink `xml:"atom:link"`
	Title         string   `xml:"title"`
	Description   string   `xml:"description"`
	Link          string   `xml:"link"`
	Language      string   `xml:"language"`
	LastBuildDate string   `xml:"lastBuildDate"`
	Items         []Item   `xml:"item"`
}

type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

type Item struct {
	Title       CDATA     `xml:"title"`
	Description CDATA     `xml:"description"`
	Link        string    `xml:"link"`
	Guid        string    `xml:"guid"`
	PubDate     string    `xml:"pubDate"`
	Enclosure   Enclosure `xml:"enclosure"`
}

type CDATA struct {
	Value string `xml:",cdata"`
}

type Enclosure struct {
	URL    string `xml:"url,attr"`
	Type   string `xml:"type,attr"`
	Length string `xml:"length,attr"`
}

func filterRSSFeed(feedURL, titlePattern string) ([]*gofeed.Item, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		return nil, err
	}

	pattern, err := regexp.Compile(titlePattern)
	if err != nil {
		return nil, err
	}

	var filteredItems []*gofeed.Item
	for _, item := range feed.Items {
		if pattern.MatchString(item.Title) {
			filteredItems = append(filteredItems, item)
		}
	}

	return filteredItems, nil
}
func buildFeed(filteredItems []*gofeed.Item, feedURL string) RSS {
	// Create RSS structure
	rss := RSS{
		Version: "2.0",
		DC:      "http://purl.org/dc/elements/1.1/",
		Content: "http://purl.org/rss/1.0/modules/content/",
		Atom:    "http://www.w3.org/2005/Atom",
		Channel: Channel{
			AtomLink: AtomLink{
				Href: feedURL,
				Rel:  "self",
				Type: "application/rss+xml",
			},
			Title:         "Filtered RSS Feed",
			Description:   "Filtered RSS feed based on the title pattern",
			Link:          feedURL,
			Language:      "en-us",
			LastBuildDate: "Mon, 17 Feb 2025 17:12:23 +0100",
		},
	}

	for _, item := range filteredItems {
		rss.Channel.Items = append(rss.Channel.Items, Item{
			Title:       CDATA{Value: item.Title},
			Description: CDATA{Value: item.Description},
			Link:        item.Link,
			Guid:        item.GUID,
			PubDate:     item.Published,
			Enclosure: Enclosure{
				URL:    item.Enclosures[0].URL,
				Type:   item.Enclosures[0].Type,
				Length: "10000",
			},
		})
	}

	return rss

}
