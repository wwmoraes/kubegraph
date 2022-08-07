package scheme

import (
	"bytes"
	"fmt"
	"go/format"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/tools/go/packages"

	_ "embed"
)

var (
	//go:embed scheme.go.gotmpl
	templateText string

	flags = struct {
		importURL      string
		outputFileName string
		importName     string
		tags           []string
	}{}

	command = &cobra.Command{
		Use: "scheme [flags] [directory]",

		Short:        "Generates scheme load glue code",
		Long:         "generates a scheme import and runtime composition to enable decoding extra resource types",
		SilenceUsage: true,
		PreRunE:      preScheme,
		RunE:         scheme,
	}
)

type TemplateData struct {
	PackageName string
	ImportName  string
	ImportURL   string
}

func preScheme(cmd *cobra.Command, args []string) error {
	if len(flags.importURL) == 0 {
		return fmt.Errorf("no import URL set")
	}

	if len(flags.importName) == 0 {
		parts := strings.Split(flags.importURL, "/")
		flags.importName = parts[len(parts)-1]
	}

	return nil
}

func scheme(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		args = []string{"."}
	}

	dir, err := filepath.Abs(filepath.Dir(args[0]))
	if err != nil {
		return fmt.Errorf("resolving path: %s", err)
	}

	if len(flags.outputFileName) == 0 {
		baseName := fmt.Sprintf("%s_scheme.go", strings.ToLower(flags.importName))
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

	schemeTemplate, err := template.New("scheme").Parse(templateText)
	if err != nil {
		return fmt.Errorf("reading template: %s", err)
	}

	var buf bytes.Buffer
	err = schemeTemplate.Execute(&buf, &TemplateData{
		ImportName:  flags.importName,
		ImportURL:   flags.importURL,
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
	command.Flags().StringVarP(&flags.outputFileName, "output", "o", "", "output file name (default srcdir/<type>_extern.go)")
	command.Flags().StringSliceVar(&flags.tags, "tags", []string{}, "comma-separated list of build tags to apply")
}

func Command() *cobra.Command {
	return command
}
