package cato

import (
	"testing"

	"github.com/cs3org/cato/resources"
)

func TestMarkdown(t *testing.T) {

	rootPath := "examples/"
	conf := &resources.CatoConfig{
		Driver: "markdown",
		DriverConfig: map[string]map[string]interface{}{
			"markdown": map[string]interface{}{
				"ReferenceBase": "https://github.com/cs3org/cato/tree/master/examples",
			},
		},
	}

	if _, err := GenerateDocumentation(rootPath, conf); err != nil {
		t.Errorf("GenerateDocumentation(): %w", err)
	}
}
