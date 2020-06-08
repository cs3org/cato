package reva

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/cs3org/cato/resources"
	"github.com/cs3org/cato/writer"
	"github.com/cs3org/cato/writer/drivers/registry"
	"github.com/mitchellh/mapstructure"
)

const (
	mdFile = "_index.md"

	configDefaultTemplate = "{{`{{%`}} dir name=\"{{ .Config.FieldName}}\" type=\"{{ .Config.DataType}}\" default={{ .Config.DefaultValue}} {{`%}}`}}\n" +
		"{{ .Config.Description}}\n" +
		"{{`{{< highlight toml >}}`}}\n" +
		"[{{ .TomlPath}}]\n" +
		"{{ .Config.FieldName}} = {{ .EscapedDefaultValue}}\n" +
		"{{`{{< /highlight >}}`}}\n" +
		"{{`{{% /dir %}}`}}\n"

	configPointerTemplate = "{{`{{%`}} dir name=\"{{ .Config.FieldName}}\" type=\"{{ .Config.DataType}}\" default=\"{{ .Config.DefaultValue}}\" {{`%}}`}}\n" +
		"{{ .Config.Description}}\n" +
		"{{`{{< highlight toml >}}`}}\n" +
		"[{{ .TomlPath}}]\n" +
		"{{ .EscapedDefaultValue}}\n" +
		"{{`{{< /highlight >}}`}}\n" +
		"{{`{{% /dir %}}`}}\n"

	headerTemplate = "---\n" +
		"title: \"{{ .Name}}\"\n" +
		"linkTitle: \"{{ .Name}}\"\n" +
		"weight: 10\n" +
		"description: >\n" +
		"  Configuration for the {{ .Name}} service\n" +
		"---"
)

func init() {
	registry.Register("reva", New)
}

type mgr struct {
	c *config
}

type config struct {
	DocPaths map[string]string
}

type templateParameters struct {
	Config              *resources.DocumentationInfo
	TomlPath            string
	EscapedDefaultValue string
}

func parseConfig(m map[string]interface{}) (*config, error) {
	c := &config{}
	if err := mapstructure.Decode(m, c); err != nil {
		return nil, err
	}
	return c, nil
}

func createMDFiles(root, mdDir string) error {
	th, err := template.New("revaHeader").Parse(headerTemplate)
	if err != nil {
		return err
	}

	err = os.MkdirAll(mdDir, 0700)
	if err != nil {
		return err
	}

	for strings.HasPrefix(mdDir, root) {
		docFile := path.Join(mdDir, mdFile)
		_, err := os.Stat(docFile)

		if err != nil {
			if os.IsNotExist(err) {
				f, err := os.Create(docFile)
				if err != nil {
					return err
				}
				defer f.Close()

				svc := struct {
					Name string
				}{
					Name: path.Base(mdDir),
				}
				b := bytes.Buffer{}
				err = th.Execute(&b, svc)
				if err != nil {
					return err
				}

				_, err = f.WriteString(b.String())
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
		mdDir = path.Dir(mdDir)
	}

	return nil
}

func New(m map[string]interface{}) (writer.ConfigWriter, error) {

	conf, err := parseConfig(m)
	if err != nil {
		return nil, fmt.Errorf("error parsing conf: %w", err)
	}

	mgr := &mgr{
		c: conf,
	}
	return mgr, nil
}

func (m mgr) WriteConfigs(configs map[string][]*resources.DocumentationInfo, filePath, rootPath string) error {

	td, err := template.New("revaDefault").Parse(configDefaultTemplate)
	if err != nil {
		return err
	}
	tp, err := template.New("revaPointer").Parse(configPointerTemplate)
	if err != nil {
		return err
	}

	docFileSuffix, err := filepath.Rel(rootPath, path.Dir(filePath))
	if err != nil {
		return err
	}

	var match string
	for k := range m.c.DocPaths {
		if strings.HasPrefix(docFileSuffix, k) && len(k) > len(match) {
			match = k
		}
	}

	configName, err := filepath.Rel(match, docFileSuffix)
	if err != nil {
		return err
	}

	docsRoot := path.Join(rootPath, m.c.DocPaths[match])
	mdDir := path.Join(docsRoot, configName)
	docFile := path.Join(mdDir, mdFile)

	err = createMDFiles(docsRoot, mdDir)
	if err != nil {
		return err
	}

	fi, err := os.Open(docFile)
	if err != nil {
		return err
	}
	defer fi.Close()

	mdCount := 0
	lines := []string{}

	scanner := bufio.NewScanner(fi)
	for scanner.Scan() {
		currLine := scanner.Text()
		lines = append(lines, currLine)
		if strings.TrimSpace(currLine) == "---" {
			mdCount = mdCount + 1
		}
		if mdCount == 2 {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	lines = append(lines, "")

	for s, fields := range configs {
		lines = append(lines, fmt.Sprintf("# _struct: %s_\n", s))

		for _, f := range fields {
			var escapedDefaultValue, tomlPath string
			var isPointer bool
			if strings.HasPrefix(f.DefaultValue, "url:") {
				decodedVal := strings.TrimPrefix(f.DefaultValue, "url:")
				escapedDefaultValue = fmt.Sprintf(`"[%s]({{< ref "%s" >}})"`, decodedVal, decodedVal)
				f.DefaultValue = decodedVal
				tomlPath = strings.ReplaceAll(configName, "/", ".") + "." + f.FieldName
				isPointer = true
			} else {
				escapedDefaultValue = f.DefaultValue
				tomlPath = strings.ReplaceAll(configName, "/", ".")
			}
			params := templateParameters{
				Config:              f,
				TomlPath:            tomlPath,
				EscapedDefaultValue: escapedDefaultValue,
			}

			b := bytes.Buffer{}
			if isPointer {
				err = tp.Execute(&b, params)
				if err != nil {
					return err
				}
			} else {
				err = td.Execute(&b, params)
				if err != nil {
					return err
				}
			}
			lines = append(lines, b.String())
		}
	}

	fo, err := os.Create(docFile)
	if err != nil {
		return err
	}
	defer fo.Close()

	w := bufio.NewWriter(fo)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()

}
