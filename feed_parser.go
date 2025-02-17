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
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
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
		Channel: Channel{
			Title:       "Filtered RSS Feed",
			Link:        feedURL,
			Description: "Filtered RSS feed based on the title pattern",
		},
	}

	for _, item := range filteredItems {
		rss.Channel.Items = append(rss.Channel.Items, Item{
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
		})
	}

	return rss
}
