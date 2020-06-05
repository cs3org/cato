package cato

import (
	"testing"

	"github.com/cs3org/cato/resources"
	_ "github.com/cs3org/cato/writer/drivers/reva"
)

func TestCato(t *testing.T) {
	rootPath := "/home/reva/"
	conf := &resources.CatoConfig{
		Driver: "reva",
		DriverConfig: map[string]map[string]interface{}{
			"reva": map[string]interface{}{
				"DocPaths": map[string]string{
					"internal/": "docs/content/en/docs/config/",
					"pkg/":      "docs/content/en/docs/config/packages/",
				},
			},
		},
	}
	if err := GenerateDocumentation(rootPath, conf); err != nil {
		t.Errorf("GenerateDocumentation(): %w", err)
	}
}
