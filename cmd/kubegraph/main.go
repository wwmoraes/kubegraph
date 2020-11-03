package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/goccy/go-graphviz"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/wwmoraes/kubegraph/icons"
	"github.com/wwmoraes/kubegraph/internal/loader"
)

var rootCmd = &cobra.Command{
	Use:     "kubegraph [file]",
	Short:   "Kubernetes resource graph generator",
	Long:    "generates a graph of kubernetes resources and their dependencies/relations",
	Args:    cobra.ExactArgs(1),
	PreRunE: preRun,
	RunE:    run,
}

var rootFlags = struct {
	outputPath string
}{}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&rootFlags.outputPath, "output-path", "o", "", "the output path of the graph dot and svg")
}

func preRun(cmd *cobra.Command, args []string) error {
	sourceFileName := args[0]

	// check if file is, uh, a file
	fileInfo, err := os.Stat(sourceFileName)
	if err != nil {
		return err
	}
	if !fileInfo.Mode().IsRegular() {
		return errors.Errorf("%s is not a valid file", sourceFileName)
	}

	if rootFlags.outputPath == "" {
		return nil
	}

	// ensure the output path exists
	if err := os.MkdirAll(rootFlags.outputPath, 0755); err != nil {
		return err
	}

	// restore icon assets
	log.Println("restoring assets...")
	icons.RestoreAssets(rootFlags.outputPath, "icons")

	return nil
}

func run(cmd *cobra.Command, args []string) error {
	sourceFileName := args[0]

	// parse file
	instance, err := loader.FromYAML(sourceFileName)
	if err != nil {
		return err
	}

	baseFileName := path.Base(strings.TrimSuffix(sourceFileName, path.Ext(sourceFileName)))
	log.Println("generating dot graph...")
	if err := instance.Render(path.Join(rootFlags.outputPath, fmt.Sprintf("%s.%s", baseFileName, "dot")), graphviz.XDOT); err != nil {
		return err
	}
	log.Println("generating svg graph...")
	if err := instance.Render(path.Join(rootFlags.outputPath, fmt.Sprintf("%s.%s", baseFileName, "svg")), graphviz.SVG); err != nil {
		return err
	}

	return nil
}
