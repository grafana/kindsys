package codegen

import (
	"fmt"
	"testing"
)

func TestTSResourceJenny_NoParams(t *testing.T) {
	test := NewGenTest(t, GenTestConfig{
		OutputDir: "testdata/codegen/output/folder_TSResourceJenny_NoParams",
	})

	test.RunOneToOneFromModule(
		"testdata/codegen/schemas/folder",
		LatestJenny("root", TSResourceJenny{
			ImportMapper: func(importStmt string) (string, error) {
				if importStmt == "github.com/grafana/kindsys" {
					return "", nil
				}

				return "", fmt.Errorf("import '%s' not expected", importStmt)
			},
		}),
	)
}
