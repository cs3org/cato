package reva

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/cs3org/cato/resources"
	"github.com/cs3org/cato/writer"
	"github.com/cs3org/cato/writer/drivers/registry"
	"github.com/mitchellh/mapstructure"
)

const (
	mdFile = "_index.md"

	configDefaultTemplate = "{{`{{%`}} dir name=\"{{ .Config.FieldName}}\" type=\"{{ .Config.DataType}}\" default=\"{{ .Config.DefaultValue}}\" {{`%}}`}}\n" +
		"{{ .Config.Description}}\n" +
		"{{`{{< highlight toml >}}`}}\n" +
		"[{{ .TomlPath}}]\n" +
		"{{ .Config.FieldName}} = \"{{ .EscapedDefaultValue}}\"\n" +
		"{{`{{< /highlight >}}`}}\n" +
		"{{`{{% /dir %}}`}}\n"

	configPointerTemplate = "{{`{{%`}} dir name=\"{{ .Config.FieldName}}\" type=\"{{ .Config.DataType}}\" default=\"{{ .Config.DefaultValue}}\" {{`%}}`}}\n" +
		"{{ .Config.Description}}\n" +
		"{{`{{< highlight toml >}}`}}\n" +
		"[{{ .TomlPath}}]\n" +
		"{{ .EscapedDefaultValue}}\n" +
		"{{`{{< /highlight >}}`}}\n" +
		"{{`{{% /dir %}}`}}\n"
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

	td, err := template.New("catoDefault").Parse(configDefaultTemplate)
	if err != nil {
		return err
	}
	tp, err := template.New("catoPointer").Parse(configPointerTemplate)
	if err != nil {
		return err
	}

	docFileSuffix := strings.TrimPrefix(path.Dir(filePath), rootPath)

	var match string
	for k := range m.c.DocPaths {
		if strings.HasPrefix(docFileSuffix, k) && len(k) > len(match) {
			match = k
		}
	}

	configName := strings.TrimPrefix(docFileSuffix, match)
	docFile := path.Join(rootPath, m.c.DocPaths[match], configName, mdFile)

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
		lines = append(lines, fmt.Sprintf("struct %s\n", s))

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
