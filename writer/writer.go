package writer

import "github.com/cs3org/cato/resources"

type ConfigWriter interface {
	WriteConfigs(configs map[string][]*resources.FieldInfo, filePath, rootPath string) error
}
