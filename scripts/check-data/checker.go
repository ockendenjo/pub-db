package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/ockendenjo/pub-db/types"
)

var (
	changingBeers = regexp.MustCompile(`serves (\d+) changing`)
	regBeers      = regexp.MustCompile(`(\d+) regular`)
)

const caskStr = "https://camra.org.uk/images/beer-containers/cask-ale-venues.png"

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

	updaters := []updater{
		updateNumBeers,
		updateCaskStatus,
	}

	hasChanges := false
	for _, updateFn := range updaters {
		changed, err := updateFn(pub, urlStr, b)
		if err != nil {
			return err
		}
		hasChanges = hasChanges || changed
	}

	if pub.RealAles > pub.NumBeers {
		pub.RealAles = pub.NumBeers
	}

	if !hasChanges {
		fmt.Printf("%s | %s - no changes\n", urlStr, pub.Name)
	}

	return nil
}

func updateCaskStatus(pub *types.Pub, urlStr string, b []byte) (bool, error) {
	hasCask := bytes.Contains(b, []byte(caskStr))
	hasChanges := false
	if hasCask != pub.HasRealAle {
		fmt.Printf("%s | %s: hasCask %t->%t\n", urlStr, pub.Name, pub.HasRealAle, hasCask)
		pub.HasRealAle = hasCask
		hasChanges = true
	}
	return hasChanges, nil
}

func updateNumBeers(pub *types.Pub, urlStr string, b []byte) (bool, error) {
	count := 0

	regMatch := regBeers.FindStringSubmatch(string(b))
	if len(regMatch) > 1 {
		regCount, err := strconv.Atoi(regMatch[1])
		if err != nil {
			return false, err
		}
		count += regCount
	}

	chgMatch := changingBeers.FindStringSubmatch(string(b))
	if len(chgMatch) > 1 {
		chgCount, err := strconv.Atoi(chgMatch[1])
		if err != nil {
			return false, err
		}
		count += chgCount
	}

	hasChanges := false
	if count != pub.NumBeers {
		fmt.Printf("%s | %s: numBeers %d->%d\n", urlStr, pub.Name, pub.NumBeers, count)
		pub.NumBeers = count
		hasChanges = true
	}
	return hasChanges, nil
}

type updater func(pub *types.Pub, urlStr string, b []byte) (bool, error)
