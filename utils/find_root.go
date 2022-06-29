// Copyright © 2022 LMU Munich Media Informatics Group. All rights reserved.
// Created by Changkun Ou <https://changkun.de>.
//
// Use of this source code is governed by a GNU GPLv3 license that
// can be found in the LICENSE file.

// This script is written for collect all lab study session id.
//
// To use this script:
//
// $ go run find_root.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type doc struct {
	UserId      string   `json:"userId"`
	IDs         []string `json:"ids"`
	Optimal     float64  `json:"optimal"`
	Root        string   `json:"root"`
	Description string
	Ranking     []struct {
		Score string `json:"id"`
		Title string `json:"title"`
		Tasks []struct {
			Title string `json:"title"`
			Id    string `json:"id"`
			Desc  string `json:"description"`
		} `json:"tasks"`
	} `json:"ranking"`
}

type Evaluation []doc

func (e Evaluation) GetAllRating(uid string) []string {
	uEval := []doc{}
	for i := range e {
		if e[i].UserId == uid {
			uEval = append(uEval, e[i])
		}
	}

	ratings := []string{}
	for i := 0; i < len(uEval); i++ {
		eval := uEval[i]
		ratings = append(ratings, eval.Root)
	}

	return ratings
}

func main() {
	var e Evaluation

	comparisons, err := os.ReadFile("../dataset/metadata/lab/seq.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(comparisons, &e)
	if err != nil {
		panic(err)
	}

	allUid := map[string]int{}
	for i := range e {
		allUid[e[i].UserId]++
	}

	f, err := os.Create("../dataset/metadata/lab/all.txt")
	if err != nil {
		panic(err)
	}
	for uid := range allUid {
		all := map[string]int{}
		for _, root := range e.GetAllRating(uid) {
			all[root] = 0
		}
		for root := range all {
			fmt.Fprintf(f, "%v\n", root)
		}
	}
	f.Close()
}
