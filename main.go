package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
	"time"
)

type FeedProgress struct {
	URL          string `json:"url"`
	FetchedUntil int64  `json:"last_fetch_time"` // Unix timestamp
}

func main() {
	log.SetLevel(log.InfoLevel)
	app := &cli.App{
		Name:  "FeedFlux (ff)",
		Usage: "FeedFlux is a lightweight tool developed in Go that parses various feeds such as RSS and Atom into a unified JSON format, with the ability to record and resume fetching progress.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "record",
				Aliases:  []string{"r"},
				Value:    "",
				Usage:    "to record fetching progress, specify a progress directory",
				Required: false,
			},

			&cli.BoolFlag{
				Name:     "continue",
				Aliases:  []string{"c"},
				Usage:    "to continue fetching from a previous record, record directory must be specified",
				Required: false,
			},

			&cli.IntFlag{
				Name:     "timeout",
				Usage:    "timeout in seconds for fetching a feed",
				Value:    10,
				Required: false,
			},
		},

		Action: func(c *cli.Context) error {
			recordDirPath := c.String("record")
			continueFetch := c.Bool("continue")
			if continueFetch && recordDirPath == "" {
				log.Fatal("record directory must be specified when continuing fetching")
			}

			urls := c.Args().Slice()
			if len(urls) == 0 {
				log.Fatal("no feed url specified")
			}
			results := make(chan *gofeed.Item)

			for _, url := range urls {
				go fetchFeed(url, results, recordDirPath, continueFetch, c.Int("timeout"))
			}

			for item := range results {
				jsonFeed, _ := json.Marshal(item)
				fmt.Println(string(jsonFeed))
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func urlToFileName(url string) string {
	fileName := sha256.Sum256([]byte(url))
	return fmt.Sprintf("%x.json", fileName)
}

func readRecordFile(url string, recordDirPath string) (*FeedProgress, error) {
	fileName := urlToFileName(url)
	recordFilePath := fmt.Sprintf("%s/%s", recordDirPath, fileName)
	recordFile, err := os.Open(recordFilePath)
	// If record file does not exist, return nil with error
	if os.IsNotExist(err) {
		return nil, err
	}
	// If other error occurs, return nil with error
	if err != nil {
		return nil, err
	}
	defer recordFile.Close()
	decoder := json.NewDecoder(recordFile)
	var record FeedProgress
	if err := decoder.Decode(&record); err != nil {
		return nil, err
	}
	return &record, nil
}

func fetchFeed(url string, results chan *gofeed.Item, recordDirPath string, continueFetch bool, timeoutSeconds int) {
    defer close(results)
	fp := gofeed.NewParser()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()
	feed, err := fp.ParseURLWithContext(url, ctx)
	if err != nil {
		log.Infof("error fetching %s: %s", url, err)
		return
	}
	// FetchedUntil is the largest published timestamp of all items
	record := FeedProgress{
		URL:          url,
		FetchedUntil: 0,
	}
	for _, item := range feed.Items {
		if item.PublishedParsed.Unix() > record.FetchedUntil {
			record.FetchedUntil = item.PublishedParsed.Unix()
		}
	}
	// record file path is sha256 of url
	fileName := sha256.Sum256([]byte(url))
	recordFilePath := fmt.Sprintf("%s/%x.json", recordDirPath, fileName)
	// log feed info
	log.WithFields(log.Fields{
		"title": feed.Title,
		"link":  feed.Link,
		"items": len(feed.Items),
	}).Infof("fetched %s", url)

	lastTimeStamp := int64(0)
	if continueFetch {
		// If continue fetching, read record file
		recordFile, err := readRecordFile(url, recordDirPath)
		if err != nil {
			log.Warnf("error reading record file %s: %s, continue fetching from scratch", recordFilePath, err)
		} else {
			log.Infof("continue fetching %s from %s", url, time.Unix(recordFile.FetchedUntil, 0))
			lastTimeStamp = recordFile.FetchedUntil
		}
	}

	for _, item := range feed.Items {
		// If continue fetching, skip items that are fetched before
		if continueFetch && item.PublishedParsed.Unix() <= lastTimeStamp {
			continue
		}
		results <- item
	}

	if recordDirPath != "" {
		// Create record directory if not exists
		if _, err := os.Stat(recordDirPath); os.IsNotExist(err) {
			if err := os.Mkdir(recordDirPath, 0755); err != nil {
				log.Fatal(err)
			}
		}
		// Write record file
		recordFile, err := os.OpenFile(recordFilePath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer recordFile.Close()
		encoder := json.NewEncoder(recordFile)
		if err := encoder.Encode(record); err != nil {
			log.Fatal(err)
		}
	}
}
