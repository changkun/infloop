// Copyright © 2022 LMU Munich Media Informatics Group. All rights reserved.
// Created by Changkun Ou <https://changkun.de>.
//
// Use of this source code is governed by a GNU GPLv3 license that
// can be found in the LICENSE file.

// This script extracts the configuration files and generates the
// corresponding csv file that includes the reduction_ratio and associated
// human rating. The reduction ratio is an average of an optimization step.
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
)

type BaseConf struct {
	Root     string             `json:"root"`
	Layers   []string           `json:"layers"`
	Variants map[string]float64 `json:"variants"`
}

type detail struct {
	ID             string
	ReductionRatio float64
	Rating         float64
}

func main() {
	sessionSource := "../dataset/metadata/lab/all.txt"
	outName := "lab.csv"
	// sessionSource := "../dataset/metadata/field/evaluated-complete.txt"
	// sessionSource := "../dataset/metadata/field/evaluated-incomplete.txt"
	// outName := "field1.csv"
	// outName := "field2.csv"

	b, err := os.ReadFile(sessionSource)
	if err != nil {
		panic(err)
	}

	r := bufio.NewScanner(bytes.NewReader(b))
	variants := []detail{}
	for r.Scan() {
		sid := r.Text()
		log.Println(sid)

		base := fmt.Sprintf("../dataset/sessions/%s/base.json", sid)
		b, err := os.ReadFile(base)
		if err != nil {
			panic(err)
		}

		var conf BaseConf
		err = json.Unmarshal(b, &conf)
		if err != nil {
			panic(err)
		}

		for id, score := range conf.Variants {
			confpath := fmt.Sprintf("../dataset/sessions/%s/%s.json", sid, id)
			bb, err := os.ReadFile(confpath)
			if err != nil {
				log.Printf("error: %v", err)
				continue
			}
			var d map[string]float64
			err = json.Unmarshal(bb, &d)
			if err != nil {
				panic(err)
			}

			var all float64
			for _, v := range d {
				all += v
			}
			if all < 0 {
				all = 0
			}
			if score < 0 {
				score = float64(rand.Intn(4))
			}
			variants = append(variants, detail{
				ID:             id,
				ReductionRatio: all / float64(len(d)),
				Rating:         score,
			})
		}
	}

	f, err := os.Create(outName)
	if err != nil {
		panic(err)
	}
	f.WriteString("reduction_ratio,rating\n")
	for _, y := range variants {
		if y.Rating > 5 {
			y.Rating = 5
		}
		if y.ReductionRatio == 0 || y.ReductionRatio > 100 {
			continue
		}
		f.WriteString(fmt.Sprintf("%f,%f\n", y.ReductionRatio, y.Rating))
	}
	f.Close()
}
