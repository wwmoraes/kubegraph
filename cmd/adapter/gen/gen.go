package gen

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/wwmoraes/kubegraph/icons"
	"golang.org/x/tools/go/packages"

	_ "embed"
)

var (
	//go:embed adapter.go.gotmpl
	templateText string

	flags = struct {
		importURL      string
		typeName       string
		outputFileName string
		iconName       string
		importName     string
		tags           []string
	}{}

	command = &cobra.Command{
		Use: "gen [flags] [directory]",

		Short:        "Generates kubegraph adapter glue code",
		Long:         "generates a struct that implements the adapter interface for a kubernetes resource, and injects its registration at runtime",
		SilenceUsage: true,
		PreRunE:      preGen,
		RunE:         gen,
	}
)

type TemplateData struct {
	PackageName string
	ImportName  string
	ImportURL   string
	TypeName    string
	IconPath    string
}

func preGen(cmd *cobra.Command, args []string) error {
	if len(flags.typeName) == 0 {
		return fmt.Errorf("no type set\n")
	}

	if len(flags.importURL) == 0 {
		return fmt.Errorf("no import URL set")
	}

	if len(strings.Split(flags.typeName, ".")) > 1 {
		return fmt.Errorf("type must not contain a package prefix")
	}

	if len(flags.importName) == 0 {
		parts := strings.Split(flags.importURL, "/")
		flags.importName = parts[len(parts)-1]
	}

	return nil
}

func gen(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		args = []string{"."}
	}

	dir, err := filepath.Abs(filepath.Dir(args[0]))
	if err != nil {
		return fmt.Errorf("resolving path: %s", err)
	}

	iconPath := fmt.Sprintf("icons/%s.svg", flags.iconName)
	if _, err = icons.Asset(iconPath); err != nil {
		return fmt.Errorf("icon not registered: %w", err)
	}

	if len(flags.outputFileName) == 0 {
		baseName := fmt.Sprintf("%s_adapter.go", strings.ToLower(flags.typeName))
		flags.outputFileName = strings.ToLower(baseName)
	}

	packageConfig := &packages.Config{
		Mode:       packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax,
		Tests:      false,
		BuildFlags: []string{fmt.Sprintf("-tags=%s", strings.Join(flags.tags, " "))},
	}

	loadedPackages, err := packages.Load(packageConfig, args...)
	if err != nil {
		return fmt.Errorf("loading package: %s", err)
	}

	if len(loadedPackages) != 1 {
		return fmt.Errorf("error: %d packages found", len(loadedPackages))
	}

	outputName := filepath.Join(dir, flags.outputFileName)

	adapterTemplate, err := template.New("adapter").Parse(templateText)
	if err != nil {
		return fmt.Errorf("reading template: %s", err)
	}

	var buf bytes.Buffer
	err = adapterTemplate.Execute(&buf, &TemplateData{
		ImportName:  flags.importName,
		ImportURL:   flags.importURL,
		TypeName:    flags.typeName,
		IconPath:    iconPath,
		PackageName: loadedPackages[0].Name,
	})
	if err != nil {
		return fmt.Errorf("executing template: %s", err)
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("internal error: invalid Go generated: %s", err)
	}

	err = ioutil.WriteFile(outputName, src, 0644)
	if err != nil {
		return fmt.Errorf("writing output: %s", err)
	}

	cmd.Printf("%s: wrote %s\n", loadedPackages[0].PkgPath, outputName)

	return nil
}

func init() {
	command.Flags().StringVarP(&flags.importURL, "import", "i", "", "import required for the target k8s resource type")
	command.Flags().StringVarP(&flags.importName, "name", "n", "", "import name to be used as the import namespace (default: last component of the import URL)")
	command.Flags().StringVarP(&flags.typeName, "type", "t", "", "target kubernetes resource type")
	command.Flags().StringVarP(&flags.outputFileName, "output", "o", "", "output file name (default srcdir/<type>_adapter.go)")
	command.Flags().StringVar(&flags.iconName, "icon", "unknown", "resource SVG icon name; must be added to the icons package separately")
	command.Flags().StringSliceVar(&flags.tags, "tags", []string{}, "comma-separated list of build tags to apply")
}

func Command() *cobra.Command {
	return command
}
