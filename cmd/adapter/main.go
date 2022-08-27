package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wwmoraes/kubegraph/cmd/adapter/dep"
	"github.com/wwmoraes/kubegraph/cmd/adapter/gen"
	"github.com/wwmoraes/kubegraph/cmd/adapter/scheme"
)

var (
	rootCmd = &cobra.Command{
		Use:   "adapter [flags] [command]",
		Short: "Manages generated adapter glue code",
		Long:  "helper tool to create adapter trees and to generate glue code for any kubernetes resource to be supported by kubegraph",
	}
)

func init() {
	rootCmd.AddCommand(
		gen.Command(),
		dep.Command(),
		scheme.Command(),
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "adapter %s failed\n", strings.Join(os.Args[1:], " "))
		os.Exit(1)
	}
}
