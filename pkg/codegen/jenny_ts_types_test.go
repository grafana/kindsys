package codegen

import (
	"fmt"
	"testing"
)

func TestTSTypesJenny_NoParams(t *testing.T) {
	test := NewGenTest(t, GenTestConfig{
		OutputDir: "testdata/codegen/output/folder_TSTypesJenny_NoParams",
	})

	test.RunOneToOneFromModule(
		"testdata/codegen/schemas/folder",
		LatestJenny("root", TSTypesJenny{
			ImportMapper: func(importStmt string) (string, error) {
				if importStmt == "github.com/grafana/kindsys" {
					return "", nil
				}

				return "", fmt.Errorf("import '%s' not expected", importStmt)
			},
		}),
	)
}
