package exporter

import "github.com/cs3org/cato/resources"

type ConfigExporter interface {
	ExportConfigs(configs map[string][]*resources.FieldInfo, filePath, rootPath string) error
}
