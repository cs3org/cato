package markdown

import (
	"github.com/cs3org/cato/resources"
	"github.com/cs3org/cato/writer"
	"github.com/cs3org/cato/writer/drivers/registry"
)

func init() {
	registry.Register("markdown", New)
}

type mgr struct {
	c *config
}

type config struct {
}

func New(m map[string]interface{}) (writer.ConfigWriter, error) {
	return &mgr{
		c: nil,
	}, nil
}

func (m mgr) WriteConfigs(configs map[string][]*resources.FieldInfo, filePath, rootPath string) error {
	return nil
}
