// Copyright © 2022 The poly.red Authors. All rights reserved.
// The use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"

	"changkun.de/x/infloop/tools/cmd"
	"changkun.de/x/infloop/tools/polyreduce-sdk-go"
	"github.com/spf13/cobra"
)

func main() {
	log.SetPrefix("infloop: ")

	var rootCmd = &cobra.Command{
		Use:   "polyred",
		Short: "A polygon reduction service",
		Long: fmt.Sprintf(`The command line tool of polygon reduction service.

Version:     %s`, polyreduce.ClientVersion),
	}
	rootCmd.AddCommand(&cobra.Command{
		Use:   "ping",
		Short: "ping polyred service",
		Run:   cmd.Ping,
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "upload [path_to_model]",
		Short: "Upload .fbx model to polyred service",
		Args:  cobra.ExactArgs(1),
		Run:   cmd.Upload,
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "config [id] [mesh_name] [target_reduction_ratio]",
		Short: "Config the simplification target",
		Args:  cobra.ExactArgs(3),
		Run:   cmd.Config,
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "run [id]",
		Short: "Trigger polygon reduction to specific model",
		Args:  cobra.ExactArgs(1),
		Run:   cmd.Run,
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "download [id] [path_to_save]",
		Short: "Download simplified model from polyred service",
		Args:  cobra.ExactArgs(2),
		Run:   cmd.Download,
	})
	rootCmd.Execute()
}
