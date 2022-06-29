// Copyright © 2022 The poly.red Authors. All rights reserved.
// The use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package cmd

import (
	"context"
	"log"
	"strconv"

	"changkun.de/x/infloop/tools/polyreduce-sdk-go"
	"github.com/spf13/cobra"
)

func Ping(cmd *cobra.Command, args []string) {
	c := polyreduce.NewClient()
	o, err := c.Ping(context.Background())
	if err != nil {
		log.Fatalf("failed to ping polyreduce service: %v", err)
	}

	log.Println("OK")
	log.Printf("%+#v", o)
}

func Upload(cmd *cobra.Command, args []string) {
	mp := args[0]

	c := polyreduce.NewClient()
	o, err := c.PolyredUpload(context.Background(), &polyreduce.PolyredUploadInput{
		ModelPath: mp,
	})
	if err != nil {
		log.Fatalf("failed to upload: %v", err)
	}

	log.Println(o.Message)
	log.Println(o.ModelId)
}
func Config(cmd *cobra.Command, args []string) {
	id := args[0]
	name := args[1]
	ratio, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		log.Fatalf("cannot parse reduction ratio: %v", ratio)
	}

	c := polyreduce.NewClient()
	err = c.PolyredConfig(context.Background(), &polyreduce.PolyredConfigInput{
		ModelID:        id,
		ReductionRatio: map[string]float64{name: ratio},
	})
	if err != nil {
		log.Fatalf("failed to config the reduction task: %v", err)
	}

	log.Println("configuration was successful.")
}
func Run(cmd *cobra.Command, args []string) {
	id := args[0]

	c := polyreduce.NewClient()
	err := c.PolyredRun(context.Background(), &polyreduce.PolyredRunInput{ModelID: id})
	if err != nil {
		log.Fatalf("failed to run: %v", err)
	}

	log.Println("simplification is complete.")
}
func Download(cmd *cobra.Command, args []string) {
	id := args[0]
	sp := args[1]

	c := polyreduce.NewClient()
	err := c.PolyredDownload(context.Background(), &polyreduce.DownloadInput{
		ModelID: id,
		Path:    sp,
	})
	if err != nil {
		log.Fatalf("failed to download: %v", err)
	}

	log.Printf("model is saved to: %s", sp)
}
