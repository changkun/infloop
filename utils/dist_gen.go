// Copyright © 2022 LMU Munich Media Informatics Group. All rights reserved.
// Created by Changkun Ou <https://changkun.de>.
//
// Use of this source code is governed by a GNU GPLv3 license that
// can be found in the LICENSE file.

// This script implements a distribution generator for overviewing all
// possible distributions of ratings for 4 different variants at a time.
//
// It is helpful to identify what might be an ideal distribution at a time.
// Since each line only represents one possibility, hence for a N times
// sequential decision problem it has much more possibilities to construct
// as a series of distributions.
package main

import (
	"fmt"
	"os"
)

func main() {
	all := [][]int{}
	for i := 0; i < 6; i++ {
		for j := 0; j < 6; j++ {
			for k := 0; k < 6; k++ {
				for l := 0; l < 6; l++ {
					all = append(all, []int{i, j, k, l})
				}
			}
		}
	}

	for _, r := range all {
		fmt.Fprintf(os.Stdout, "%v\n", r)
	}
}
