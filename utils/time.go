// Copyright © 2022 LMU Munich Medieninformatik. All rights reserved.
// Created by Changkun Ou <https://changkun.de>.
//
// Use of this source code is governed by a GNU GPLv3 license that
// can be found in the LICENSE file.

// This script generates a time.csv file that includes all raw OS creation
// time of all files appeared in the dataset.
//
// Due to historical design reasons, some UUIDs are using version 4 which
// did not include the information when did the file being created. Hence,
// the time.csv file is helpful if we would like to understand whole
// behavior of an evaluation sequence.
//
// It is not necessary to run this code. Instead, use dataset/metadata/time.csv
// directly.
package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type row struct {
	path string
	time time.Time
}

func main() {
	allTime := []row{}
	filepath.Walk("../dataset/sessions", func(path string, info fs.FileInfo, err error) error {
		if strings.Contains(path, "gitkeep") {
			return nil
		}

		allTime = append(allTime, row{path, info.ModTime()})
		return nil
	})

	f, err := os.Create("../dataset/metadata/time.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("path,time\n")
	for _, t := range allTime {
		f.WriteString(fmt.Sprintf("%s,%v\n", strings.TrimPrefix(strings.TrimPrefix(t.path, "../dataset/sessions"), "/"), t.time.UTC()))
	}
}
