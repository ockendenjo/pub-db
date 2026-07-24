package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ockendenjo/pub-db/types"
)

func main() {
	var id int
	flag.IntVar(&id, "id", 0, "Pub ID")
	flag.Parse()

	errLog := log.New(os.Stderr, "", 0)
	err := filepath.WalkDir("pubs", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".json") {
			if ferr := processFile(path, id, errLog); ferr != nil {
				errLog.Printf("%s: %s\n", path, ferr.Error())
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func processFile(file string, id int, errLog *log.Logger) error {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	var pf types.PubsFile
	err = json.NewDecoder(f).Decode(&pf)
	if err != nil {
		return err
	}

	httpClient := &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
		},
	}
	c := &checker{httpClient: httpClient}

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
		return err
	}
	if _, err = f.Seek(0, io.SeekStart); err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err = enc.Encode(pf); err != nil {
		return err
	}
	return nil
}
