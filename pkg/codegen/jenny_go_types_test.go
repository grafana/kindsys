package codegen

import (
	"testing"
)

func TestGoTypesJenny_NoParams(t *testing.T) {
	test := NewGenTest(t, GenTestConfig{
		OutputDir: "testdata/codegen/output/folder_GoTypesJenny_NoParams",
	})

	test.RunOneToOneFromModule(
		"testdata/codegen/schemas/folder",
		LatestJenny("parent", GoTypesJenny{}),
	)
}
