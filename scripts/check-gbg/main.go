package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/ockendenjo/pub-db/types"
)

func main() {
	f, err := os.Open("pubs/pubs.json")
	if err != nil {
		panic(err)
	}
	var pf types.PubsFile
	err = json.NewDecoder(f).Decode(&pf)
	if err != nil {
		panic(err)
	}

	httpClient := &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
		},
	}
	c := &checker{httpClient: httpClient}

	errLog := log.New(os.Stderr, "", 0)

	for _, pub := range pf.Pubs {
		err = c.checkPub(pub)
		if err != nil {
			errLog.Printf("%s %s\n", pub.Name, err.Error())
		}
		time.Sleep(1 * time.Second)
	}
}

type checker struct {
	httpClient *http.Client
}

func (c *checker) checkPub(pub *types.Pub) error {
	if pub.GoodBeerID == 0 {
		return nil
	}

	gbgStr := strconv.Itoa(pub.GoodBeerID)

	urlStr, err := url.JoinPath("https://goodbeerguide.org.uk/pub/", gbgStr, "/show")
	if err != nil {
		return err
	}

	res, err := c.httpClient.Get(urlStr)
	if err != nil {
		return err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		fmt.Printf("%s: HTTP status %d\n", pub.Name, res.StatusCode)
		return nil
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	count := 0

	regMatch := regBeers.FindStringSubmatch(string(b))
	if len(regMatch) > 1 {
		regCount, err := strconv.Atoi(regMatch[1])
		if err != nil {
			return err
		}
		count += regCount
	}

	chgMatch := changingBeers.FindStringSubmatch(string(b))
	if len(chgMatch) > 1 {
		chgCount, err := strconv.Atoi(chgMatch[1])
		if err != nil {
			return err
		}
		count += chgCount
	}

	if count != pub.RealAles {
		fmt.Printf("%s | %s: realAles %d->%d\n", urlStr, pub.Name, pub.RealAles, count)
		return nil
	}

	//Any more checks?
	return nil
}

var regBeers = regexp.MustCompile(`(\d+) regular beer`)
var changingBeers = regexp.MustCompile(`(\d+) changing beer`)
