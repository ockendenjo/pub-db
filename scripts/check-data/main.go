package main

import (
	"bytes"
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
	f, err := os.OpenFile("pubs/pubs.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = f.Close()
	}()

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
		if false {
			err = c.checkPub(pub)
			if err != nil {
				errLog.Printf("%s %s\n", pub.Name, err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}

	if err = f.Truncate(0); err != nil {
		panic(err)
	}
	if _, err = f.Seek(0, io.SeekStart); err != nil {
		panic(err)
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err = enc.Encode(pf); err != nil {
		panic(err)
	}
}

type checker struct {
	httpClient *http.Client
}

func (c *checker) checkPub(pub *types.Pub) error {
	if pub.GoodBeerID == nil {
		return nil
	}

	idStr := strconv.Itoa(pub.CamraID)

	urlStr, err := url.JoinPath("https://camra.org.uk/pubs/", idStr)
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

	if err = updateNumBeers(pub, urlStr, b); err != nil {
		return err
	}
	if pub.RealAles > pub.NumBeers {
		pub.RealAles = pub.NumBeers
	}
	if err = updateCaskStatus(pub, urlStr, b); err != nil {
		return err
	}

	return nil
}

func updateCaskStatus(pub *types.Pub, urlStr string, b []byte) error {
	hasCask := bytes.ContainsAny(b, caskStr)
	if hasCask != pub.HasRealAle {
		fmt.Printf("%s | %s: hasCask %t->%t\n", urlStr, pub.Name, pub.HasRealAle, hasCask)
		pub.HasRealAle = hasCask
	}
	return nil
}

func updateNumBeers(pub *types.Pub, urlStr string, b []byte) error {
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

	if count != pub.NumBeers {
		fmt.Printf("%s | %s: numBeers %d->%d\n", urlStr, pub.Name, pub.NumBeers, count)
		pub.NumBeers = count
	}
	return nil
}

var changingBeers = regexp.MustCompile(`serves (\d+) changing`)
var regBeers = regexp.MustCompile(`(\d+) regular`)

const caskStr = "https://camra.org.uk/images/beer-containers/cask-ale-venues.png"
