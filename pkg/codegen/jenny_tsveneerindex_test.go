package codegen

import (
	"fmt"
	"testing"
)

func TestTSVeneerIndexJenny_Basic(t *testing.T) {
	test := NewGenTest(t, GenTestConfig{
		OutputDir: "testdata/codegen/output/folder_TSVeneerIndex_Basic",
	})

	test.RunManyToOneFromModule(
		"testdata/codegen/schemas/folder",
		TSVeneerIndexJenny("dir", func(importStmt string) (string, error) {
			if importStmt == "github.com/grafana/kindsys" {
				return "", nil
			}

			return "", fmt.Errorf("import '%s' not expected", importStmt)
		}),
	)
}
