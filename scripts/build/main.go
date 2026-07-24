package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/ockendenjo/pub-db/types"
)

func main() {
	allPubs := make([]*types.Pub, 0)

	err := filepath.WalkDir("pubs", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".json") {
			pubs, err := processFile(path)
			if err != nil {
				panic(err)
			}
			allPubs = append(allPubs, pubs...)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll("output", 0755)
	if err != nil {
		panic(err)
	}

	pf := types.PubsFileWithoutSchema{Pubs: allPubs}
	b, err := json.Marshal(pf)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("output/pubs.json", b, 0644)
	if err != nil {
		panic(err)
	}
}

func processFile(path string) ([]*types.Pub, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var pf types.PubsFileWithoutSchema
	err = json.Unmarshal(b, &pf)
	if err != nil {
		return nil, err
	}
	return pf.Pubs, nil
}
