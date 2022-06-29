// Copyright © 2022 LMU Munich Media Informatics Group. All rights reserved.
// Created by Changkun Ou <https://changkun.de>.
//
// Use of this source code is governed by a GNU GPLv3 license that
// can be found in the LICENSE file.

// This script checks the number of iterations that the user participated.
package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

func main() {
	field()
	lab()
}

func field() {
	bcomp, err := os.ReadFile("../dataset/metadata/field/evaluated-complete.txt")
	if err != nil {
		panic(err)
	}

	bincomp, err := os.ReadFile("../dataset/metadata/field/evaluated-incomplete.txt")
	if err != nil {
		panic(err)
	}

	seqComplete := strings.Split(string(bcomp), "\n")
	seqIncomplete := strings.Split(string(bincomp), "\n")
	fmt.Println("total: ", len(seqComplete)+len(seqIncomplete))

	iters := []float64{}
	for _, id := range seqComplete {
		iters = append(iters, countIterations(id))
	}
	for _, id := range seqIncomplete {
		iters = append(iters, countIterations(id))
	}

	fmt.Println("field study:")
	fmt.Println("[min, max]: ", min(iters), max(iters))
	fmt.Println("mu: ", mean(iters))
	fmt.Println("sigma: ", std(iters))
}

func lab() {
	b, err := os.ReadFile("../dataset/metadata/lab/all.txt")
	if err != nil {
		panic(err)
	}
	seq := strings.Split(string(b), "\n")
	iters := []float64{}
	for _, id := range seq {
		iters = append(iters, countIterations(id))
	}
	fmt.Println("lab study:")
	fmt.Println("[min, max]: ", min(iters), max(iters))
	fmt.Println("mu: ", mean(iters))
	fmt.Println("sigma: ", std(iters))
}

func countIterations(id string) float64 {
	dir := "../dataset/sessions/" + id
	entries, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	// exlucde base model and config, 4 models per iteration.
	return float64((len(entries) - 2) / 8)
}
func min(a []float64) float64 {
	min := math.MaxFloat64
	for _, v := range a {
		if v < min {
			min = v
		}
	}
	return min
}
func max(a []float64) float64 {
	max := -math.MaxFloat64
	for _, v := range a {
		if v > max {
			max = v
		}
	}
	return max
}

func mean(a []float64) float64 {
	sum := 0.0
	for _, v := range a {
		sum += v
	}
	return sum / float64(len(a))
}

func std(a []float64) float64 {
	mu := mean(a)

	sum := 0.0
	for _, v := range a {
		sum += (v - mu) * (v - mu)
	}
	return math.Sqrt(sum / float64(len(a)))
}
