package cato

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/cs3org/cato/resources"
	"github.com/cs3org/cato/writer"
	_ "github.com/cs3org/cato/writer/drivers/loader"
	"github.com/cs3org/cato/writer/drivers/registry"
)

type structInfo struct {
	StructDef  *ast.StructType
	StructName string
}

var namedTags = []string{"xml", "mapstructure", "json"}

func listGoFiles(rootPath string) ([]string, error) {
	goFileRegex, _ := regexp.Compile(`^.+\.go$`)
	fileList := []string{}

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if goFileRegex.MatchString(info.Name()) {
			fileList = append(fileList, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return fileList, nil
}

func getLineInitialPositions(filePath string) ([]int, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var charCount int
	initPositions := []int{0}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		charCount = charCount + len(scanner.Text())
		initPositions = append(initPositions, charCount)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return initPositions, nil
}

func getLineNumber(lineNos []int, pos int) (int, error) {
	for i, n := range lineNos {
		if pos <= n {
			return i, nil
		}
	}

	return -1, fmt.Errorf("position exceeds total characters in the file")
}

func parseStruct(structDef *ast.StructType, catoTag, filePath string, fset *token.FileSet, lineNos []int) ([]*resources.FieldInfo, error) {
	configs := []*resources.FieldInfo{}

	for _, field := range structDef.Fields.List {
		if field.Tag != nil {

			tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))
			configTag := tag.Get(catoTag)

			if configTag != "" {
				// get field.Type as string
				var typeNameBuf bytes.Buffer
				err := printer.Fprint(&typeNameBuf, fset, field.Type)
				if err != nil {
					return nil, fmt.Errorf("error decoding struct field name: %w", err)
				}

				var fieldName string
				for _, namedTag := range namedTags {
					if t := tag.Get(namedTag); t != "" {
						fieldName = strings.Split(t, ",")[0]
					}
				}
				if fieldName == "" {
					fieldName = field.Names[0].Name
				}

				var desc string
				if field.Doc != nil {
					comments := []string{}
					for _, c := range field.Doc.List {
						c.Text = strings.ReplaceAll(c.Text, "//", "")
						c.Text = strings.ReplaceAll(c.Text, "/*", "")
						c.Text = strings.ReplaceAll(c.Text, "*/", "")
						c.Text = strings.Join(strings.Fields(c.Text), " ")
						comments = append(comments, c.Text)
					}
					desc = strings.Join(comments, " ")
				}

				var defaultVal string

				switch splitVals := strings.Split(configTag, ";"); len(splitVals) {
				case 1:
					defaultVal = splitVals[0]
				case 2:
					defaultVal = splitVals[0]
					desc = splitVals[1]
				case 3:
					fieldName = splitVals[0]
					defaultVal = splitVals[1]
					desc = splitVals[2]
				}

				if typeNameBuf.String() == "string" {
					defaultVal = fmt.Sprintf("\"%s\"", defaultVal)
				}

				lineNumber, err := getLineNumber(lineNos, int(field.Pos()))
				if err != nil {
					return nil, err
				}

				configs = append(configs, &resources.FieldInfo{
					FieldName:    fieldName,
					DefaultValue: defaultVal,
					Description:  desc,
					DataType:     typeNameBuf.String(),
					LineNumber:   lineNumber,
				})
			}
		}
	}

	return configs, nil
}

func getConfigsToDocument(filePath, catoTag string) (map[string][]*resources.FieldInfo, error) {
	fset := token.NewFileSet()
	fileTree, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	lineNos, err := getLineInitialPositions(filePath)
	if err != nil {
		return nil, err
	}

	structList := []*structInfo{}
	configs := map[string][]*resources.FieldInfo{}

	ast.Inspect(fileTree, func(node ast.Node) bool {
		spec, ok := node.(*ast.TypeSpec)
		if !ok {
			return true
		}
		s, ok := spec.Type.(*ast.StructType)
		if !ok {
			return true
		}
		structList = append(structList, &structInfo{s, spec.Name.Name})
		return false
	})

	for _, s := range structList {
		c, err := parseStruct(s.StructDef, catoTag, filePath, fset, lineNos)
		if err != nil {
			return nil, err
		}
		if len(c) > 0 {
			configs[s.StructName] = c
		}
	}
	return configs, nil
}

func getDriver(c *resources.CatoConfig) (writer.ConfigWriter, error) {
	if f, ok := registry.NewFuncs[c.Driver]; ok {
		return f(c.DriverConfig[c.Driver])
	}
	return nil, fmt.Errorf("driver not found: %s", c.Driver)
}

func GenerateDocumentation(rootPath string, conf *resources.CatoConfig) error {

	if rootPath == "" {
		return fmt.Errorf("cato: root path can't be empty")
	}

	if conf.CustomTag == "" {
		conf.CustomTag = "docs"
	}

	fileList, err := listGoFiles(rootPath)
	if err != nil {
		return fmt.Errorf("cato: error listing root path: %w", err)
	}

	writerDriver, err := getDriver(conf)
	if err != nil {
		return fmt.Errorf("cato: error getting driver: %w", err)
	}

	for _, file := range fileList {

		configs, err := getConfigsToDocument(file, conf.CustomTag)
		if err != nil {
			return fmt.Errorf("cato: error parsing go file: %w", err)
		}

		if len(configs) > 0 {
			err = writerDriver.WriteConfigs(configs, file, rootPath)
			if err != nil {
				return fmt.Errorf("cato: error writing documentation: %w", err)
			}
		}
	}
	return nil
}
