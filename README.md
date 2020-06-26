# Cato

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go](https://github.com/cs3org/cato/workflows/Go/badge.svg)](https://github.com/cs3org/cato/actions)

Cato is an automated documentation generation library for Go Projects. Through the use of custom tags for struct fields, Cato can extract information such as
- data type of the fields
- default values
- description through comments or tags
- online references to point to

After this extraction, this data can be exported through multiple interfaces, such as HTML and markdown.

The motivation for creating this library was to save developers from the trouble of describing configuration details separately for the end-users. With Cato, it can be generated on the fly by adding minimal info through custom tags.


## Installation

`go get github.com/cs3org/cato`


## Usage

Some examples can be found in [cato_html_test.go](cato_html_test.go) and [cato_markdown_test.go](cato_markdown_test.go) which act on [filesystem.go](examples/filesystem.go). The general usage is:

```go
import (
	"github.com/cs3org/cato"
	"github.com/cs3org/cato/resources"
)

func main() {
	rootPath := "examples/"
	conf := &resources.CatoConfig{
		Driver: "markdown",
		DriverConfig: map[string]map[string]interface{}{
			"markdown": map[string]interface{}{
				"ReferenceBase": "https://github.com/cs3org/cato/tree/master/examples",
			},
		},
	}

	if _, err := cato.GenerateDocumentation(rootPath, conf); err != nil {
		log.Error(err)
	}
}
```

We've integrated Cato with [Reva](https://github.com/cs3org/reva/) using a make [rule](https://github.com/cs3org/reva/blob/master/tools/generate-documentation/main.go) with a custom driver, where it's being used in [production](https://reva.link/docs/config/grpc/services/storageprovider/).

## Workflow

### Extraction

Cato works by generating the syntax tree for the go files using the [parser](https://golang.org/pkg/go/parser/) and [ast](https://golang.org/pkg/go/ast/) packages. It then inspects the tree to find structs with fields possessing a custom tag (the default is `docs`), and extracts details about those into the `FieldInfo` struct.

A maximum of three values, separated by semicolons can be defined in these custom tags. The expected order of these values is:
1. The name of the field as it should appear in the docs. If this is not specified, it looks for a few commonly used tags, namely `xml`, `mapstructure` and `json`, to pick up the field name from. If none of these are found, it uses the actual name of the field.
2. The default value which is used for that particular field if it is not specified by the user. This makes it really convenient for end users reading the documentation to understand the configuration specifics.
3. A description of the field. If no description is provided, Cato reads the comments provided with the field.

As an example, the `FileSystem` struct defined below lists the various ways in which tags can be defined.

```go
type FileSystem struct {

	// docs:"field_name;default_value;description"
	CacheDirectory     string   `docs:"cache_directory;/var/tmp/;Path of cache directory"`

	// docs:"default_value;description" with json tag
	EnableLogging      bool     `json:"enable_logging" docs:"false;Whether to enable logging"`

	// docs:"default_value;description"
	AvailableChecksums []string `docs:"[adler, rabin];The list of checksums provided by the file system"`

	// docs:"default_value" with comment
	// Configs for various metadata drivers
	DriverConfig map[string]map[string]interface{} `docs:"{json:{encoding: UTF8}, xml:{encoding: ASCII}}"`
}
```

### Exporters

This extracted information can be exported through multiple interfaces including markdown and HTML. If paths for the the documentation files are specified, the files are created there, otherwise they are exported in the same directory as the go file. If a reference address is provided, a pointer to the line numbers in a remotely hosted repo is also added for each of the fields.


## License

Cato is distributed under the [Apache 2.0 license](https://github.com/cs3org/cato/blob/master/LICENSE).
