package dep

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/tools/go/packages"

	_ "embed"
)

var (
	//go:embed extern.go.gotmpl
	templateText string

	flags = struct {
		importURL      string
		typeName       string
		outputFileName string
		importName     string
		tags           []string
		prefixed       bool
	}{}

	command = &cobra.Command{
		Use: "dep [flags] [directory]",

		Short:        "Generates external resource type adapters glue code",
		Long:         "generates type aliases and a getter method to retrieve an adapter for a dependency resource type",
		SilenceUsage: true,
		PreRunE:      preDep,
		RunE:         dep,
	}
)

type TemplateData struct {
	PackageName string
	ImportName  string
	ImportURL   string
	TypeName    string
	Prefix      string
}

func preDep(cmd *cobra.Command, args []string) error {
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

func dep(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		args = []string{"."}
	}

	dir, err := filepath.Abs(filepath.Dir(args[0]))
	if err != nil {
		return fmt.Errorf("resolving path: %s", err)
	}

	if len(flags.outputFileName) == 0 {
		baseName := fmt.Sprintf("%s%s_extern.go", strings.ToLower(flags.importName), strings.ToLower(flags.typeName))
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

	externTemplate, err := template.New("extern").Parse(templateText)
	if err != nil {
		return fmt.Errorf("reading template: %s", err)
	}

	prefix := ""
	if flags.prefixed {
		prefix = cases.Title(language.Und, cases.NoLower).String(flags.importName)
	}

	var buf bytes.Buffer
	err = externTemplate.Execute(&buf, &TemplateData{
		ImportName:  flags.importName,
		ImportURL:   flags.importURL,
		TypeName:    flags.typeName,
		PackageName: loadedPackages[0].Name,
		Prefix:      prefix,
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
	command.Flags().StringVarP(&flags.outputFileName, "output", "o", "", "output file name (default srcdir/<type>_extern.go)")
	command.Flags().StringSliceVar(&flags.tags, "tags", []string{}, "comma-separated list of build tags to apply")
	command.Flags().BoolVar(&flags.prefixed, "prefixed", false, "prepend symbol names with import name (useful for versioned dependencies)")
}

func Command() *cobra.Command {
	return command
}
