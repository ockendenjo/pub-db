package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ockendenjo/pub-db/types"
)

func main() {
	var id int
	flag.IntVar(&id, "id", 0, "Pub ID")
	flag.Parse()

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
		if id > 0 && pub.CamraID != id {
			continue
		}

		err = c.checkPub(pub)
		if err != nil {
			errLog.Printf("%s %s\n", pub.Name, err.Error())
		}
		time.Sleep(1 * time.Second)
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
