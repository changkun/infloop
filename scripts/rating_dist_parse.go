// Copyright © 2022 LMU Munich Media Informatics Group. All rights reserved.
// Created by Changkun Ou <https://changkun.de>.
//
// Use of this source code is governed by a GNU GPLv3 license that
// can be found in the LICENSE file.

// This script parses ratings distributions of cherry-picked model IDs from
// the field study. All data are generated into data/ratingdist folder.
//
// Usage:
//
// $ go run rating_dist_parse.go
package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"
)

// Cherry-picked from the field study.
var ids = []string{
	"cd0db28f-9277-49a1-a0e1-ab4d684fe572",
	"086d8b7a-fa86-4b80-93b7-1a40fed9bc80",
	"ea144125-cf0f-4ef4-9b3d-53737f8f33d0",
	"dacdd183-4dd0-11ec-86eb-a85e4557a9b6",
	"bfbb6dda-9c5c-44e4-ad36-f451eca40fc5",
	"921f50e8-f995-486b-bb78-f82558cd55da",
	"0c2b7579-72d8-483a-a90c-f9ec5ff591ca",
	"75ea0ac8-5584-4c89-b4f3-d5972ced14ff",
	"47a89f48-2633-4e3e-b6af-139e6cd6c513",
	"48f9c537-43f4-4a22-a7ff-4ce8d9e48d4e",
	"a0864fcd-520e-428f-8f5b-d09b254c36cb",
	"33cc17f4-2cc9-40be-99e2-91f29d64986f",
}

func main() {
	all := parseAllTime()

	for _, id := range ids {
		path := fmt.Sprintf("../dataset/sessions/%s/base.json", id)

		b, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}

		conf := &base{}
		err = json.Unmarshal(b, conf)
		if err != nil {
			panic(err)
		}

		f, err := os.Create(fmt.Sprintf("./data/ratingdist/%s.csv", id))
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(f, "root,model,unixTime,iteration,rating,\n")
		sorted := []row{}

		for modelId, rating := range conf.Variants {
			sorted = append(sorted, row{Root: id, Model: modelId, UnixTime: all[modelId].Unix(), Rating: int(rating)})
		}
		sort.SliceStable(sorted, func(i, j int) bool {
			return sorted[i].UnixTime < sorted[j].UnixTime
		})

		iteration := 0
		for _, r := range sorted {
			if r.Rating < 0 {
				continue
			}
			fmt.Fprintf(f, "%v,%v,%v,%v,%v\n", r.Root, r.Model, r.UnixTime, iteration/4, r.Rating)
			iteration++
		}
		f.Close()
	}

}

func parseAllTime() map[string]time.Time {
	b, err := os.ReadFile("../dataset/metadata/time.csv")
	if err != nil {
		panic(err)
	}
	all := map[string]time.Time{}
	o := NewScanner(strings.NewReader(string(b)))
	for o.Scan() {
		if !strings.Contains(o.Text("path"), ".fbx") {
			continue
		}

		splits := strings.Split(o.Text("path"), "/")
		if len(splits) < 2 {
			continue
		}

		tt := strings.TrimSuffix(o.Text("time"), " +0000 UTC")
		t, err := time.Parse("2006-01-02 15:04:05.99999", tt)
		if err != nil {
			panic(err)
		}

		all[strings.TrimSuffix(splits[1], ".fbx")] = t
	}
	return all
}

type base struct {
	Root     string             `json:"root"`
	Variants map[string]float64 `json:"variants"`
}

type row struct {
	Root, Model string
	UnixTime    int64
	Rating      int
}

type Scanner struct {
	Reader *csv.Reader
	Head   map[string]int
	Row    []string
}

func NewScanner(o io.Reader) Scanner {
	csv_o := csv.NewReader(o)
	a, e := csv_o.Read()
	if e != nil {
		return Scanner{}
	}
	m := map[string]int{}
	for n, s := range a {
		m[s] = n
	}
	return Scanner{Reader: csv_o, Head: m}
}

func (o *Scanner) Scan() bool {
	a, e := o.Reader.Read()
	o.Row = a
	return e == nil
}

func (o Scanner) Text(s string) string {
	return o.Row[o.Head[s]]
}
