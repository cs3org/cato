package cato

import (
	"testing"

	"github.com/cs3org/cato/resources"
)

func TestHTML(t *testing.T) {

	rootPath := "examples/"
	conf := &resources.CatoConfig{
		Driver: "html",
	}

	if _, err := GenerateDocumentation(rootPath, conf); err != nil {
		t.Errorf("GenerateDocumentation(): %w", err)
	}
}
