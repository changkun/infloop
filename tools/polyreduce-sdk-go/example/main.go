// Copyright © 2022 LMU Munich Media Informatics Group. All rights reserved.
// Created by Changkun Ou <https://changkun.de>.
//
// Use of this source code is governed by a GNU GPLv3 license that
// can be found in the LICENSE file.

// This script generates the dataset/reduces.

package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"changkun.de/x/infloop/tools/polyreduce-sdk-go"
)

var client = polyreduce.NewClient()

func simplify(modelPath, modelID string, ratio map[string]float64) {
	// 1. upload config
	err := client.PolyredConfig(context.Background(), &polyreduce.PolyredConfigInput{
		ModelID:        modelID,
		ReductionRatio: ratio,
	})
	if err != nil {
		panic(err)
	}

	// 2. run simplification
	err = client.PolyredRun(context.Background(), &polyreduce.PolyredRunInput{
		ModelID: modelID,
	})
	if err != nil {
		panic(err)
	}

	// 3. download models back
	p := 0.0
	for _, v := range ratio {
		p += v
	}
	p /= float64(len(ratio))

	err = client.PolyredDownload(context.Background(), &polyreduce.DownloadInput{
		ModelID: modelID,
		Path:    fmt.Sprintf(modelPath+"_%v.fbx", p),
	})
	if err != nil {
		panic(err)
	}
}

type UploadOutput struct {
	ID  string `json:"id"`
	Msg string `json:"msg"`
}

type ConfigInput struct {
	Percent map[string]float64 `json:"percent"`
}

func process(model string) {
	modelPath := model + "_0.obj"
	b, err := os.ReadFile(modelPath)
	if err != nil {
		panic(err)
	}
	layers := []string{}
	all := strings.Split(string(b), "\n")
	for i := range all {
		if !strings.HasPrefix(all[i], "o ") {
			continue
		}

		layers = append(layers, strings.TrimPrefix(all[i], "o "))
	}

	modelPath = model + "_0.fbx"

	o, err := client.PolyredUpload(context.Background(), &polyreduce.PolyredUploadInput{
		ModelPath: modelPath,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(o.ModelId)
	configs := []*ConfigInput{}
	for i := 1; i < 100; i++ {
		m := map[string]float64{}
		for _, l := range layers {
			m[l] = float64(i)
		}
		configs = append(configs, &ConfigInput{Percent: m})
	}

	for _, conf := range configs {
		simplify(model, o.ModelId, conf.Percent)
	}
}

func main() {
	var models = []string{
		"../dataset/reduces/teapot/teapot",
		"../dataset/reduces/cow/cow",
		"../dataset/reduces/monkey/monkey",
		"../dataset/reduces/pumpkin/pumpkin",
		"../dataset/reduces/rose/rose",
	}

	for _, model := range models {
		process(model)
	}
}
